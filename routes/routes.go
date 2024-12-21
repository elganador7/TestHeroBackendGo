package routes

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, agent *agent.Agent, cfg *config.Config) {
	setupAuthRoutes(router, db)
	setupUserAnswerRoutes(router, db)
	setupQuestionRoutes(router, db)
	setupQuestionAnswerRoutes(router, db)
	setupQueryRoutes(router, db, agent)
	setupTestTopicDataRoutes(router, db)
}

func SetupTestRoutes(router *gin.Engine, db *gorm.DB) {
	setupAuthRoutes(router, db)
	setupUserAnswerRoutesTest(router, db)
	setupQuestionRoutesTest(router, db)
	setupQuestionAnswerRoutesTest(router, db)
	setupTestTopicDataRoutesTest(router, db)
}
