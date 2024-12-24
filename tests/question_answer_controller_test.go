package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForQuestionAnswer() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.QuestionAnswer{})
	return db
}

func setupRouterWithQuestionAnswerController(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	controller := controllers.NewQuestionAnswerController(db)
	router.GET("/api/answers/question/:questionId", controller.GetAnswerByQuestionID)
	return router
}

func TestGetAnswerByQuestionID(t *testing.T) {
	db := setupTestDBForQuestionAnswer()
	router := setupRouterWithQuestionAnswerController(db)

	// Seed a question-answer record
	answer := models.QuestionAnswer{
		QuestionID:    "123",
		CorrectAnswer: "The correct answer is 42",
	}
	db.Create(&answer)

	// Test case: Valid question ID
	req := httptest.NewRequest(http.MethodGet, "/api/answers/question/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.QuestionAnswer
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, answer.QuestionID, response.QuestionID)
	assert.Equal(t, answer.CorrectAnswer, response.CorrectAnswer)

	// Test case: Invalid question ID
	req = httptest.NewRequest(http.MethodGet, "/api/answers/question/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Equal(t, "Answer not found", errorResponse["error"])
}
