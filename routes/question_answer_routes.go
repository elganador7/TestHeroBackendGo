package routes

import (
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQuestionAnswerRoutes(router *gin.Engine, db *gorm.DB, isTest bool) {
	questionAnswerCtrl := controllers.NewQuestionAnswerController(db)

	questionAnswerApi := router.Group("/api/questionAnswers")

	// Protected routes
	questionAnswerApi.Use(utils.GenerateHandlers(isTest)...)
	{
		questionAnswerApi.GET("/:questionId", questionAnswerCtrl.GetAnswerByQuestionID)
	}
}
