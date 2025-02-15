package controllers

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/agent/prompts"
	"TestHeroBackendGo/models"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DEFAULT_CORRECT_SCORE       = 0.5
	DEFAULT_QUESTION_DIFFICULTY = 0.5
	MINIMUM_QUESTION_DIFFICULTY = 0.05
)

type QueryController struct {
	DB    *gorm.DB
	Agent *agent.Agent
}

func NewQueryController(db *gorm.DB, agent *agent.Agent) *QueryController {
	return &QueryController{
		Agent: agent,
		DB:    db,
	}
}

func (ctrl *QueryController) GenerateNewQuestionHandler(c *gin.Context) {
	var req struct {
		TestTopicData models.TestTopicData `json:"test_topic_data"`
		UserId        string               `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := ctrl.GenerateNewQuestionWithTopicData(req.TestTopicData, rand.Float64())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

func (ctrl *QueryController) GenerateSimilarQuestionHandler(c *gin.Context) {
	questionId := c.Param("questionId")
	var originalQuestion models.Question

	// Fetch the question by ID
	if err := ctrl.DB.First(&originalQuestion, "id = ?", questionId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question"})
			return
		}
	}

	inputSchema := models.SimilarQuestionGeneratorInputSchema{
		Paragraph:    originalQuestion.Paragraph,
		QuestionText: originalQuestion.QuestionText,
	}

	systemPrompt, ok := prompts.SubjectTopicPromptMap[originalQuestion.TestTopic.TestType][originalQuestion.TestTopic.Subject]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test type or subject in original question"})
		return
	}

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateSimilarQuestion(inputSchema, systemPrompt)
	if err != nil {
		log.Printf("Error generating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate question"})
		return
	}

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Printf("Error generating answer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate answer"})
		return
	}

	optionInput := models.OptionGeneratorInputSchema{
		QuestionText:  questionResponse.QuestionText,
		Explanation:   answerResponse.Explanation,
		CorrectAnswer: answerResponse.CorrectAnswer,
	}

	optionsResponse, err := ctrl.Agent.GenerateQuestionOptions(optionInput)
	if err != nil {
		log.Printf("Error generating options: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate options"})
		return
	}

	// Save the question to the database
	question := models.Question{
		ID:            uuid.NewString(),
		QuestionText:  questionResponse.QuestionText,
		TestTopicID:   originalQuestion.TestTopicID,
		TestTopic:     originalQuestion.TestTopic,
		Options:       optionsResponse.Options,
		EstimatedTime: 60,
	}

	if err := ctrl.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question"})
		return
	}

	answer := models.QuestionAnswer{
		ID:            uuid.NewString(),
		QuestionID:    question.ID,
		CorrectAnswer: optionsResponse.CorrectOption,
		Explanation:   answerResponse.Explanation,
	}

	if err := ctrl.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// GenerateQuestion generates a question tailored to the current user's performance
func (ctrl *QueryController) GenerateRelevantQuestion(c *gin.Context) {
	var req struct {
		TestType string `json:"test_type"`
		Subject  string `json:"subject"`
		UserId   string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if testType and subject are provided
	if req.TestType == "" || req.Subject == "" || req.UserId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "req.TestType, Subject, and UserID are required"})
		return
	}

	// Query for all topics, subtopics, and specific topics under the req.TestType and subject
	var testTopics []models.TestTopicData
	if err := ctrl.DB.Where("test_type = ? AND subject = ?", req.TestType, req.Subject).Find(&testTopics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test topics"})
		return
	}

	// Query the user performance summary
	var userPerformance []models.UserPerformanceSummary
	if err := ctrl.DB.Where("user_id = ? AND test_type = ? AND subject = ?", req.UserId, req.TestType, req.Subject).Find(&userPerformance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user performance"})
		return
	}

	// Create a map of specific topics to performance
	performanceMap := make(map[string]float64)
	questionDifficultyMap := make(map[string]float64)
	for _, performance := range userPerformance {
		if performance.TotalPointsPossible == 0 {
			performanceMap[performance.SpecificTopic] = DEFAULT_CORRECT_SCORE
			questionDifficultyMap[performance.SpecificTopic] = DEFAULT_QUESTION_DIFFICULTY
		} else {
			difficulty := performance.TotalPoints / performance.TotalPointsPossible
			performanceMap[performance.SpecificTopic] = difficulty * performance.TotalPoints
			questionDifficultyMap[performance.SpecificTopic] = difficulty
		}
	}

	// Create a list of topics with their weights
	var selectedTopic models.TestTopicData
	minWeight := 1000000000000000.0 // Very big number
	difficulty := 0.5

	for _, topic := range testTopics {
		// Get the correct rate for the specific topic
		score, ok := performanceMap[topic.SpecificTopic]
		if !ok {
			// If no data available for this topic, treat it as 50% correct rate
			score = DEFAULT_CORRECT_SCORE
			performanceMap[topic.SpecificTopic] = score
			questionDifficultyMap[topic.SpecificTopic] = DEFAULT_QUESTION_DIFFICULTY
		}

		// Calculate weight: lower correct rate means higher weight
		weight := score * rand.Float64()
		if weight < minWeight {
			minWeight = weight
			if questionDifficultyMap[topic.SpecificTopic] > MINIMUM_QUESTION_DIFFICULTY {
				difficulty = questionDifficultyMap[topic.SpecificTopic]
			} else {
				difficulty = MINIMUM_QUESTION_DIFFICULTY
			}
			selectedTopic = topic
		}
	}

	// Query for a question based on the selected topic
	question, err := ctrl.GenerateNewQuestionWithTopicData(selectedTopic, difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to generate question: %v", err)})
		log.Printf("Failed to generate question: %v", err)
		return
	}

	// Return the selected question
	c.JSON(http.StatusOK, question)
}

func (ctrl *QueryController) GenerateNewQuestionWithTopicData(testTopicData models.TestTopicData, difficulty float64) (models.Question, error) {
	previousQuestionModels := []models.Question{}
	previousQuestionTexts := []string{}

	if err := ctrl.DB.Where("test_topic_id = ?", testTopicData.ID).Find(&previousQuestionModels).Limit(5).Error; err != nil {
		return models.Question{}, err
	}

	for _, question := range previousQuestionModels {
		previousQuestionTexts = append(previousQuestionTexts, question.QuestionText)
	}

	inputSchema := models.NewQuestionGeneratorInputSchema{
		Topic:             testTopicData.Topic,
		Subtopic:          testTopicData.Subtopic,
		SpecificTopic:     testTopicData.SpecificTopic,
		Difficulty:        difficulty,
		PreviousQuestions: previousQuestionTexts,
	}

	systemPrompt, ok := prompts.SubjectTopicPromptMap[testTopicData.TestType][testTopicData.Subject]
	if !ok {
		return models.Question{}, fmt.Errorf("no prompt found for test type %s and subject %s", testTopicData.TestType, testTopicData.Subject)
	}

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateNewQuestion(inputSchema, systemPrompt)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		return models.Question{}, err
	}

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Printf("Error generating answer: %v", err)
		return models.Question{}, err
	}

	optionInput := models.OptionGeneratorInputSchema{
		QuestionText:  questionResponse.QuestionText,
		Explanation:   answerResponse.Explanation,
		CorrectAnswer: answerResponse.CorrectAnswer,
	}

	optionsResponse, err := ctrl.Agent.GenerateQuestionOptions(optionInput)
	if err != nil {
		log.Printf("Error generating options: %v", err)
		return models.Question{}, err
	}

	// Save the question to the database
	question := models.Question{
		ID:            uuid.NewString(),
		QuestionText:  questionResponse.QuestionText,
		Options:       optionsResponse.Options,
		EstimatedTime: 60,
		Difficulty:    inputSchema.Difficulty,
		TestTopicID:   testTopicData.ID,
		TestTopic:     testTopicData,
	}

	formattedQuestion, err := ctrl.Agent.ValidateMathJaxFormatting(question)
	if err != nil {
		return models.Question{}, err
	}

	question = formattedQuestion.TranslateQuestionOutputSchemaToQuestion(question)

	if err := ctrl.DB.Create(&question).Error; err != nil {
		return models.Question{}, err
	}

	answer := models.QuestionAnswer{
		ID:            uuid.NewString(),
		QuestionID:    question.ID,
		CorrectAnswer: optionsResponse.CorrectOption,
		Explanation:   answerResponse.Explanation,
	}

	if err := ctrl.DB.Create(&answer).Error; err != nil {
		return models.Question{}, err
	}

	return question, nil
}
