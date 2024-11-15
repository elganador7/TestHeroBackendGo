package controllers

import (
	"TestHeroBackendGo/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
