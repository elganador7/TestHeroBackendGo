package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.UserAnswer{}, &models.UserPerformanceSummary{})
	return db
}

func setupRouterWithController(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	controller := controllers.NewUserAnswerController(db)
	router.POST("/api/user_answers/", controller.CreateUserAnswer)
	router.GET("/api/user_answers/user/:userId", controller.GetUserAnswersByUser)
	router.POST("/api/user_answers/user/summary", controller.GetUserPerformanceSummary)
	return router
}

func TestCreateUserAnswer(t *testing.T) {
	db := setupTestDB()
	router := setupRouterWithController(db)

	input := models.UserAnswer{
		UserID:     "1",
		QuestionID: "2",
		TimeTaken:  30,
		Attempts:   1,
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
	assert.NotEmpty(t, createdAnswer.ID)
	assert.WithinDuration(t, time.Now(), createdAnswer.CreatedAt, time.Second)
}

func TestGetUserAnswersByUser(t *testing.T) {
	db := setupTestDB()
	router := setupRouterWithController(db)

	answer := models.UserAnswer{
		ID:         uuid.New().String(),
		UserID:     "1",
		QuestionID: "2",
		TimeTaken:  30,
		Attempts:   1,
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
	router := gin.Default()
	agent := agent.NewAgent("", db, "")
	routes.SetupRoutes(router, db, agent, true)

	db.Exec(`INSERT INTO user_performance_summary (user_id, correct_rate) VALUES
		('1', 0.75)`)

	input := map[string]string{
		"userId": "1",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/user_answers/user/summary", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.UserPerformanceSummary
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response, 1)
	assert.Equal(t, "1", response[0].UserID)
	assert.Equal(t, 0.75, response[0].CorrectRate)
}
