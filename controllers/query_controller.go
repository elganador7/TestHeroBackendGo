package controllers

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/models"
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

	log.Printf("Request body: %v", req)

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

	log.Printf("Input schema: %v", inputSchema)

	// Call the agent with the system prompt
	questionResponse, err := ctrl.Agent.GenerateNewQuestion(inputSchema)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate question"})
		return
	}

	log.Printf("Generated question: %v", questionResponse)

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Fatalf("Error generating answer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate answer"})
		return
	}

	log.Printf("Generated answer: %v", answerResponse)

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

	log.Printf("Generated options: %v", optionsResponse)

	// Save the question to the database
	question := models.Question{
		ID:            uuid.NewString(),
		QuestionText:  questionResponse.QuestionText,
		TestType:      req.TestType,
		Subject:       req.Subject,
		Topic:         req.Topic,
		Subtopic:      req.Subtopic,
		Options:       optionsResponse.Options,
		EstimatedTime: 60,
		Difficulty:    inputSchema.Difficulty,
	}

	log.Printf("Saved question: %v", question)

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

	log.Printf("Original question: %v", originalQuestion)

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

	log.Printf("Generated question: %v", questionResponse)

	answerResponse, err := ctrl.Agent.GenerateAnswer(questionResponse)
	if err != nil {
		log.Fatalf("Error generating answer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate answer"})
		return
	}

	log.Printf("Generated answer: %v", answerResponse)

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
		CorrectAnswer: optionsResponse.CorrectOption,
		Explanation:   answerResponse.Explanation,
	}

	if err := ctrl.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
		return
	}

	c.JSON(http.StatusOK, question)
}
