package controllers

import (
	"net/http"
	"strconv"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

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
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var results []struct {
		SubjectArea string
		CorrectRate float64
	}

	query := `
		SELECT subject_area, AVG(CASE WHEN is_correct THEN 1 ELSE 0 END) AS correct_rate
		FROM user_answers
		WHERE user_id = ?
		GROUP BY subject_area
	`
	if err := ctrl.DB.Raw(query, userID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate performance summary"})
		return
	}

	c.JSON(http.StatusOK, results)
}
