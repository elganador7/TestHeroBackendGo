package controllers

import (
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /questions
func GetQuestions(c *gin.Context) {
	var questions []models.Question
	database.DB.Find(&questions)
	c.JSON(http.StatusOK, questions)
}

// POST /questions
func CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&question)
	c.JSON(http.StatusOK, question)
}
