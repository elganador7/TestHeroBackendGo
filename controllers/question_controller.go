package controllers

import (
	"TestHeroBackendGo/models"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QuestionController struct {
	DB *gorm.DB
}

func NewQuestionController(db *gorm.DB) *QuestionController {
	return &QuestionController{DB: db}
}

// GET /questions
func (ctrl *QuestionController) GetQuestions(c *gin.Context) {
	var questions []models.Question
	ctrl.DB.Find(&questions)
	c.JSON(http.StatusOK, questions)
}

// POST /questions
func (ctrl *QuestionController) CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Create(&question)
	c.JSON(http.StatusOK, question)
}

func (ctrl *QuestionController) CreateQuestionWithAnswer(c *gin.Context) {
	var req struct {
		QuestionText  string            `json:"question_text"`
		TestType      string            `json:"test_type"`
		Subject       string            `json:"subject"`
		Topic         string            `json:"topic"`
		Difficulty    float64           `json:"difficulty"`
		Options       datatypes.JSONMap `json:"options"`
		EstimatedTime int               `json:"estimated_time"`
		CorrectAnswer string            `json:"correct_answer"`
		Explanation   string            `json:"explanation"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create question
	question := models.Question{
		QuestionText:  req.QuestionText,
		TestType:      req.TestType,
		Subject:       req.Subject,
		Topic:         req.Topic,
		Difficulty:    req.Difficulty,
		Options:       req.Options,
		EstimatedTime: req.EstimatedTime,
	}
	if err := ctrl.DB.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	// Create answer
	answer := models.QuestionAnswer{
		QuestionID:    question.ID,
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
	}
	if err := ctrl.DB.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Question and answer created successfully"})
}

func (ctrl *QuestionController) GetQuestionByID(c *gin.Context) {
	questionId := c.Param("id")
	var question models.Question

	// Fetch the question by ID
	if err := ctrl.DB.First(&question, "id = ?", questionId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question"})
		}
	}

	c.JSON(http.StatusOK, question)
}

func (ctrl *QuestionController) GetRandomQuestion(c *gin.Context) {
	// Database context
	var question models.Question
	var count int64

	// Count the total number of questions
	if err := ctrl.DB.Model(&models.Question{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count questions"})
		return
	}

	log.Printf("Total questions: %d", count)

	if count == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No questions available"})
		return
	}

	// Generate a random offset
	offset := rand.Intn(int(count))

	log.Printf("Random offset: %d", offset)

	// Fetch the question at the random offset
	if err := ctrl.DB.Offset(offset).Limit(1).Find(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch random question"})
		return
	}

	log.Println("Random question:", question)

	// Return the random question
	c.JSON(http.StatusOK, question)
}
