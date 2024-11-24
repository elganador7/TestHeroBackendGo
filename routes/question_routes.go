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
		questionApi.GET("/allQuestions", questionCtrl.GetQuestions)
		questionApi.GET("/:questionId", questionCtrl.GetQuestionByID)
		questionApi.POST("/create", questionCtrl.CreateQuestionWithAnswer)
		questionApi.GET("/random", questionCtrl.GetRandomQuestion)
	}
}

func setupQuestionRoutesTest(router *gin.Engine, db *gorm.DB) {
	questionCtrl := controllers.NewQuestionController(db)

	questionApi := router.Group("/api/questions")

	// Protected routes
	questionApi.Use()
	{
		questionApi.GET("/allQuestions", questionCtrl.GetQuestions)
		questionApi.GET("/:questionId", questionCtrl.GetQuestionByID)
		questionApi.POST("/create", questionCtrl.CreateQuestionWithAnswer)
		questionApi.GET("/random", questionCtrl.GetRandomQuestion)
	}
}
