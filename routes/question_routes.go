package routes

import (
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQuestionRoutes(router *gin.Engine, db *gorm.DB, isTest bool) {
	questionCtrl := controllers.NewQuestionController(db)

	questionApi := router.Group("/api/questions")

	// Protected routes
	questionApi.Use(utils.GenerateHandlers(isTest)...)
	{
		questionApi.GET("/:questionId", questionCtrl.GetQuestionByID)
		questionApi.POST("/create", questionCtrl.CreateQuestionWithAnswer)
		questionApi.GET("/random", questionCtrl.GetRandomQuestion)
	}
}
