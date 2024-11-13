package routes

import (
	"TestHeroBackendGo/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/questions", controllers.GetQuestions)
		api.POST("/questions", controllers.CreateQuestion)
	}
}
