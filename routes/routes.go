package routes

import (
	"TestHeroBackendGo/agent"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, agent *agent.Agent, stripeKey string, isTest bool) {
	setupAuthRoutes(router, db)
	setupUserAnswerRoutes(router, db, isTest)
	setupQuestionRoutes(router, db, isTest)
	setupQuestionAnswerRoutes(router, db, isTest)
	setupQueryRoutes(router, db, agent, isTest)
	setupTestTopicDataRoutes(router, db, isTest)
	setupStripeRoutes(router, stripeKey, isTest)
}
