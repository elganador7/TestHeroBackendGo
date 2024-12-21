package controllers

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/models"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Request body for OpenAI API
type OpenAIRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
	Model     string `json:"model"`
}

// OpenAIResponse structure
type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

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
		TestType string `json:"test_type"`
		Subject  string `json:"subject"`
		Topic    string `json:"topic"`
		Subtopic string `json:"subtopic"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputSchema := models.NewQuestionGeneratorInputSchema{
		TestType:   req.TestType,
		Subject:    req.Subject,
		Topic:      req.Topic,
		Subtopic:   req.Subtopic,
		Difficulty: rand.Float64(),
	}

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateNewQuestion(inputSchema)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate question"})
		return
	}

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Fatalf("Error generating answer: %v", err)
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
		log.Fatalf("Error generating options: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate options"})
		return
	}

	// Save the question to the database
	question := models.Question{
		ID:            uuid.NewString(),
		QuestionText:  questionResponse.QuestionText,
		Options:       optionsResponse.Options,
		EstimatedTime: 60,
		Difficulty:    inputSchema.Difficulty,
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

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateSimilarQuestion(inputSchema)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate question"})
		return
	}

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Fatalf("Error generating answer: %v", err)
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
		log.Fatalf("Error generating options: %v", err)
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

	testType := req.TestType
	subject := req.Subject
	userID := req.UserId

	// Check if testType and subject are provided
	if testType == "" || subject == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TestType, Subject, and UserID are required"})
		return
	}

	// Query for all topics, subtopics, and specific topics under the testType and subject
	var testTopics []models.TestTopicData
	if err := ctrl.DB.Where("test_type = ? AND subject = ?", testType, subject).Find(&testTopics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test topics"})
		return
	}

	// Query the user performance summary
	var userPerformance []models.UserPerformanceSummary
	if err := ctrl.DB.Where("user_id = ? AND test_type = ? AND subject = ?", userID, testType, subject).Find(&userPerformance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user performance"})
		return
	}

	// Create a map of specific topics to performance
	performanceMap := make(map[string]float64)
	for _, performance := range userPerformance {
		performanceMap[performance.SpecificTopic] = performance.CorrectRate
	}

	// Create a list of topics with their weights
	var weightedTopics []models.TestTopicData
	weightMap := make(map[string]float64)

	for _, topic := range testTopics {
		// Get the correct rate for the specific topic
		correctRate, ok := performanceMap[topic.SpecificTopic]
		if !ok {
			// If no data available for this topic, treat it as 50% correct rate
			correctRate = 0.5
			performanceMap[topic.SpecificTopic] = correctRate
		}

		// Calculate weight: lower correct rate means higher weight
		weight := 1.01 - (correctRate)
		weightMap[topic.SpecificTopic] = weight
		weightedTopics = append(weightedTopics, topic)
	}

	// Now, adjust each topic's weight by multiplying it by a random number between 0 and 1
	var maxWeight float64
	var selectedTopic models.TestTopicData

	log.Printf("weightMap: %v", weightMap)

	// Find the topic with the highest adjusted weight
	for _, topic := range weightedTopics {
		// Multiply weight by a random value between 0 and 1
		randomFactor := rand.Float64() // Random float between 0 and 1
		adjustedWeight := weightMap[topic.SpecificTopic] * randomFactor

		// Select the topic with the highest adjusted weight
		if adjustedWeight > maxWeight {
			maxWeight = adjustedWeight
			selectedTopic = topic
		}
	}

	// Query for a question based on the selected topic
	question, err := ctrl.GenerateNewQuestionWithTopicData(selectedTopic, performanceMap[selectedTopic.SpecificTopic])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to generate question: %v", err)})
	}

	// Return the selected question
	c.JSON(http.StatusOK, question)
}

func (ctrl *QueryController) GenerateNewQuestionWithTopicData(testTopicData models.TestTopicData, difficulty float64) (models.Question, error) {
	inputSchema := models.NewQuestionGeneratorInputSchema{
		TestType:      testTopicData.TestType,
		Subject:       testTopicData.Subject,
		Topic:         testTopicData.Topic,
		Subtopic:      testTopicData.Subtopic,
		SpecificTopic: testTopicData.SpecificTopic,
		Difficulty:    difficulty,
	}

	log.Printf("Test Topic Data: %+v", testTopicData)

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateNewQuestion(inputSchema)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		return models.Question{}, err
	}

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Fatalf("Error generating answer: %v", err)
		return models.Question{}, err
	}

	optionInput := models.OptionGeneratorInputSchema{
		QuestionText:  questionResponse.QuestionText,
		Explanation:   answerResponse.Explanation,
		CorrectAnswer: answerResponse.CorrectAnswer,
	}

	optionsResponse, err := ctrl.Agent.GenerateQuestionOptions(optionInput)
	if err != nil {
		log.Fatalf("Error generating options: %v", err)
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
