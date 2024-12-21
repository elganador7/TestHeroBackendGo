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
		userAnswersApi.POST("/submitUserAnswer", userAnswerCtrl.CreateUserAnswer)
		userAnswersApi.GET("/user/:userId", userAnswerCtrl.GetUserAnswersByUser)
		userAnswersApi.POST("/user/summary", userAnswerCtrl.GetUserPerformanceSummary)
	}
}

func setupUserAnswerRoutesTest(router *gin.Engine, db *gorm.DB) {
	userAnswerCtrl := controllers.UserAnswerController{DB: db}

	userAnswersApi := router.Group("/api/user_answers")

	userAnswersApi.Use()
	{
		userAnswersApi.POST("/", userAnswerCtrl.CreateUserAnswer)
		userAnswersApi.GET("/user/:userId", userAnswerCtrl.GetUserAnswersByUser)
		userAnswersApi.GET("/user/:userId/summary", userAnswerCtrl.GetUserPerformanceSummary)
	}
}
