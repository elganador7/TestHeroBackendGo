package routes

import (
	"TestHeroBackendGo/auth"
	"TestHeroBackendGo/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupUserAnswerRoutes(router *gin.Engine, db *gorm.DB) {
	userAnswerCtrl := controllers.UserAnswerController{DB: db}

	userAnswersApi := router.Group("/api/user_answers")

	userAnswersApi.Use(auth.JWTAuthMiddleware())
	{
		userAnswersApi.POST("/", userAnswerCtrl.CreateUserAnswer)
		userAnswersApi.GET("/user/:userId", userAnswerCtrl.GetUserAnswersByUser)
		userAnswersApi.GET("/user/:userId/summary", userAnswerCtrl.GetUserPerformanceSummary)
	}
}
