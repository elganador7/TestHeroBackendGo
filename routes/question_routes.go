package routes

import (
	"TestHeroBackendGo/auth"
	"TestHeroBackendGo/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQuestionRoutes(router *gin.Engine, db *gorm.DB) {
	questionCtrl := controllers.NewQuestionController(db)

	questionApi := router.Group("/api/questions")

	// Protected routes
	questionApi.Use(auth.JWTAuthMiddleware())
	{
		questionApi.GET("/questions", questionCtrl.GetQuestions)
		questionApi.POST("/questions", questionCtrl.CreateQuestion)
	}
}
