package controllers

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/models"
	"log"
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

// func (ctrl *QueryController) GenerateNewQuestionHandler(c *gin.Context) {
// 	var input struct {
// 		TestType string `json:"test_type"`
// 		Subject  string `json:"subject"`
// 		Topic    string `json:"topic"`
// 		Subtopic string `json:"subtopic"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Call the agent with the system prompt
// 	response, err := ctrl.Agent.GenerateQuestionWithAnswer(input.TestType, input.Subject, input.Topic, input.Subtopic)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate question"})
// 		return
// 	}

// 	// Parse the generated question
// 	var questionData map[string]interface{}
// 	if err := json.Unmarshal([]byte(response), &questionData); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid question format"})
// 		return
// 	}

// 	// Save the question to the database
// 	question := models.Question{
// 		ID:            uuid.NewString(),
// 		QuestionText:  questionData["question_text"].(string),
// 		TestType:      input.TestType,
// 		Subject:       input.Subject,
// 		Topic:         input.Topic,
// 		Subtopic:      input.Subtopic,
// 		Options:       datatypes.JSONMap(questionData["options"].(map[string]interface{})),
// 		EstimatedTime: int(questionData["estimated_time"].(float64)),
// 	}

// 	if err := ctrl.DB.Create(&question).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question"})
// 		return
// 	}

// 	answer := models.QuestionAnswer{
// 		ID:            uuid.NewString(),
// 		QuestionID:    question.ID,
// 		CorrectAnswer: questionData["correct_answer"].(string),
// 		Explanation:   questionData["explanation"].(string),
// 	}

// 	if err := ctrl.DB.Create(&answer).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, question)
// }

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

	log.Printf("Original question: %v", originalQuestion)

	inputSchema := models.QuestionGeneratorInputSchema{
		Paragraph:    originalQuestion.Paragraph,
		QuestionText: originalQuestion.QuestionText,
	}

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateQuestionWithAnswer(inputSchema)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate question"})
		return
	}

	log.Printf("Generated question: %v", questionResponse)

	optionsResponse, err := ctrl.Agent.GenerateQuestionOptions(questionResponse)
	if err != nil {
		log.Fatalf("Error generating options: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate options"})
		return
	}

	log.Printf("Generated options: %v", optionsResponse)

	// Save the question to the database
	question := models.Question{
		ID:            uuid.NewString(),
		QuestionText:  questionResponse.QuestionText,
		TestType:      originalQuestion.TestType,
		Subject:       originalQuestion.Subject,
		Topic:         originalQuestion.Topic,
		Subtopic:      originalQuestion.Subtopic,
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
		CorrectAnswer: optionsResponse.CorrectAnswer,
		Explanation:   questionResponse.CorrectAnswer,
	}

	if err := ctrl.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
		return
	}

	c.JSON(http.StatusOK, question)
}
