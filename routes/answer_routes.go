package routes

import (
	"TestHeroBackendGo/auth"
	"TestHeroBackendGo/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQuestionAnswerRoutes(router *gin.Engine, db *gorm.DB) {
	questionAnswerCtrl := controllers.NewQuestionAnswerController(db)

	questionAnswerApi := router.Group("/api/questionAnswers")

	// Protected routes
	questionAnswerApi.Use(auth.JWTAuthMiddleware())
	{
		questionAnswerApi.GET("/:questionId", questionAnswerCtrl.GetAnswerByQuestionID)
	}
}

func setupQuestionAnswerRoutesTest(router *gin.Engine, db *gorm.DB) {
	questionAnswerCtrl := controllers.NewQuestionAnswerController(db)

	// Protected routes
	questionAnswerApi := router.Group("/api/questionAnswers")

	questionAnswerApi.Use()
	{
		questionAnswerApi.GET("/:questionId", questionAnswerCtrl.GetAnswerByQuestionID)
	}
}
