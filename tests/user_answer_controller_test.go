package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.UserAnswer{})
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db)
	return router
}

func TestCreateUserAnswer(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	input := models.UserAnswer{
		UserID:      1,
		QuestionID:  2,
		AnsweredAt:  time.Now(),
		TimeTaken:   30,
		IsCorrect:   true,
		SubjectArea: "Math",
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/user_answers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdAnswer models.UserAnswer
	db.First(&createdAnswer)
	assert.Equal(t, input.UserID, createdAnswer.UserID)
}

func TestGetUserAnswersByUser(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	answer := models.UserAnswer{
		UserID:      1,
		QuestionID:  2,
		AnsweredAt:  time.Now(),
		TimeTaken:   30,
		IsCorrect:   true,
		SubjectArea: "Math",
	}
	db.Create(&answer)

	req := httptest.NewRequest(http.MethodGet, "/api/user_answers/user/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.UserAnswer
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response, 1)
	assert.Equal(t, answer.UserID, response[0].UserID)
}

func TestGetUserPerformanceSummary(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	db.Create(&models.UserAnswer{
		UserID:      1,
		QuestionID:  1,
		IsCorrect:   true,
		SubjectArea: "Math",
	})
	db.Create(&models.UserAnswer{
		UserID:      1,
		QuestionID:  2,
		IsCorrect:   false,
		SubjectArea: "Math",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/user_answers/user/1/summary", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response, 1)
	assert.Equal(t, "Math", response[0]["subject_area"])
	assert.Equal(t, 0.5, response[0]["correct_rate"])
}
