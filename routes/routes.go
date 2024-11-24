package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	setupAuthRoutes(router)
	setupUserAnswerRoutes(router, db)
	setupQuestionRoutes(router, db)
	setupQuestionAnswerRoutes(router, db)
}

func SetupTestRoutes(router *gin.Engine, db *gorm.DB) {
	setupAuthRoutes(router)
	setupUserAnswerRoutesTest(router, db)
	setupQuestionRoutesTest(router, db)
	setupQuestionAnswerRoutesTest(router, db)
}
