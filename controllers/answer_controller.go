package controllers

import (
	"TestHeroBackendGo/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionAnswerController struct {
	DB *gorm.DB
}

func NewQuestionAnswerController(db *gorm.DB) *QuestionAnswerController {
	return &QuestionAnswerController{DB: db}
}

func (ctrl *QuestionAnswerController) GetAnswerByQuestionID(c *gin.Context) {
	id := c.Param("questionId")
	var answer models.QuestionAnswer

	if err := ctrl.DB.Where("question_id = ?", id).First(&answer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	c.JSON(http.StatusOK, answer)
}
