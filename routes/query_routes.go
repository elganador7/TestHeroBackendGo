package routes

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/auth"
	"TestHeroBackendGo/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQueryRoutes(router *gin.Engine, db *gorm.DB, agent *agent.Agent) {
	queryCtrl := controllers.NewQueryController(db, agent)

	queryApi := router.Group("/api/oai_queries")

	// Protected routes
	queryApi.Use(auth.JWTAuthMiddleware())
	{
		queryApi.GET("/generate/similar/:questionId", queryCtrl.GenerateSimilarQuestionHandler)
		queryApi.POST("/generate/new", queryCtrl.GenerateNewQuestionHandler)
	}
}
