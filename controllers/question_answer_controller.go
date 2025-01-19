package controllers

import (
	"TestHeroBackendGo/models"
	"log"
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
	questionId := c.Param("questionId")
	log.Printf("questionId is %s", questionId)
	answer := models.QuestionAnswer{}

	if err := ctrl.DB.Where(models.QuestionAnswer{QuestionID: questionId}).First(&answer).Error; err != nil {
		log.Printf("Failed to get question answer: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	c.JSON(http.StatusOK, answer)
}
