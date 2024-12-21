package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"TestHeroBackendGo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAnswerController struct {
	DB *gorm.DB
}

func NewUserAnswerController(db *gorm.DB) *UserAnswerController {
	return &UserAnswerController{DB: db}
}

// CreateUserAnswer handles creating a new user answer record
func (ctrl *UserAnswerController) CreateUserAnswer(c *gin.Context) {
	var input models.UserAnswer
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Input: %v", input)

	input.ID = uuid.New().String()
	input.CreatedAt = time.Now()

	if err := ctrl.DB.Create(&input).Error; err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	log.Printf("Created record: %v", input)

	c.JSON(http.StatusCreated, input)
}

// GetUserAnswersByUser retrieves all answers for a specific user
func (ctrl *UserAnswerController) GetUserAnswersByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var answers []models.UserAnswer
	if err := ctrl.DB.Where("user_id = ?", userID).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch answers"})
		return
	}

	c.JSON(http.StatusOK, answers)
}

// GetUserPerformanceSummary retrieves a summary of the user's performance
func (ctrl *UserAnswerController) GetUserPerformanceSummary(c *gin.Context) {
	type Input struct {
		UserId string `json:"userId"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var performance []models.UserPerformanceSummary

	// Query the materialized view for the user's performance by subtopic
	if err := ctrl.DB.
		Table("user_performance_summary").
		Where("user_id = ?", input.UserId).
		Find(&performance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate performance summary"})
		return
	}

	log.Printf("Results: %v", performance)

	c.JSON(http.StatusOK, performance)
}
