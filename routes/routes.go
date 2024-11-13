package routes

import (
	"TestHeroBackendGo/auth"
	"TestHeroBackendGo/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Auth routes
	api.POST("/register", auth.Register)
	api.POST("/login", auth.Login)

	// Protected routes
	api.Use(auth.JWTAuthMiddleware())
	{
		api.GET("/questions", controllers.GetQuestions)
		api.POST("/questions", controllers.CreateQuestion)
	}
}
