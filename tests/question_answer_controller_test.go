package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBForQuestionAnswer() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(models.AllModels...)
	return db
}

func TestGetAnswerByQuestionID(t *testing.T) {
	db := setupTestDBForQuestionAnswer()
	router := gin.Default()
	agent := agent.NewAgent("", db)
	routes.SetupRoutes(router, database.DB, agent, true)

	// Seed a question-answer record
	answer := models.QuestionAnswer{
		QuestionID:    "123",
		CorrectAnswer: "The correct answer is 42",
	}
	db.Create(&answer)

	// Test case: Valid question ID
	req := httptest.NewRequest(http.MethodGet, "/api/questionAnswers/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.QuestionAnswer
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, answer.QuestionID, response.QuestionID)
	assert.Equal(t, answer.CorrectAnswer, response.CorrectAnswer)

	// Test case: Invalid question ID
	req = httptest.NewRequest(http.MethodGet, "/api/questionAnswers/999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Equal(t, "Answer not found", errorResponse["error"])
}
