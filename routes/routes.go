package routes

import (
	"TestHeroBackendGo/agent"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, agent *agent.Agent) {
	setupAuthRoutes(router)
	setupUserAnswerRoutes(router, db)
	setupQuestionRoutes(router, db)
	setupQuestionAnswerRoutes(router, db)
	setupQueryRoutes(router, db, agent)
}

func SetupTestRoutes(router *gin.Engine, db *gorm.DB) {
	setupAuthRoutes(router)
	setupUserAnswerRoutesTest(router, db)
	setupQuestionRoutesTest(router, db)
	setupQuestionAnswerRoutesTest(router, db)
}
