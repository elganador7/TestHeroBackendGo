package routes

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupQueryRoutes(router *gin.Engine, db *gorm.DB, agent *agent.Agent, isTest bool) {
	queryCtrl := controllers.NewQueryController(db, agent)

	queryApi := router.Group("/api/oai_queries")

	// Protected routes
	queryApi.Use(utils.GenerateHandlers(isTest)...)
	{
		queryApi.GET("/generate/similar/:questionId", queryCtrl.GenerateSimilarQuestionHandler)
		queryApi.POST("/generate/new", queryCtrl.GenerateNewQuestionHandler)
		queryApi.POST("/generate/relevant", queryCtrl.GenerateRelevantQuestion)
	}
}
