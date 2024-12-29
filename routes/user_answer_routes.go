package routes

import (
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupUserAnswerRoutes(router *gin.Engine, db *gorm.DB, isTest bool) {
	userAnswerCtrl := controllers.UserAnswerController{DB: db}

	userAnswersApi := router.Group("/api/user_answers")

	userAnswersApi.Use(utils.GenerateHandlers(isTest)...)
	{
		userAnswersApi.POST("/submitUserAnswer", userAnswerCtrl.CreateUserAnswer)
		userAnswersApi.GET("/user/:userId", userAnswerCtrl.GetUserAnswersByUser)
		userAnswersApi.POST("/user/summary", userAnswerCtrl.GetUserPerformanceSummary)
	}
}
