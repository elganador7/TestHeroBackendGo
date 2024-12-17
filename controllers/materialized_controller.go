package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserPerformanceSummary struct {
	UserID      int     `json:"user_id"`
	SubjectArea string  `json:"subject_area"`
	CorrectRate float64 `json:"correct_rate"`
}

type PerformanceController struct {
	DB *gorm.DB
}

func (ctrl *PerformanceController) GetUserPerformanceSummary(c *gin.Context) {
	// Parse user ID from request parameters
	userID := c.Param("userId")

	// Query the materialized view
	var results []UserPerformanceSummary
	if err := ctrl.DB.Raw(`
		SELECT user_id, subject_area, correct_rate 
		FROM user_performance_summary 
		WHERE user_id = ?`, userID).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user performance summary"})
		return
	}

	// Return results as JSON
	c.JSON(http.StatusOK, results)
}
