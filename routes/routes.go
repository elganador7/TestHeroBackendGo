package routes

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/controllers" // Import the controllers package

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Updated function signature to include paymentController
func SetupRoutes(router *gin.Engine, db *gorm.DB, agent *agent.Agent, paymentController *controllers.PaymentController, isTest bool) {
	setupAuthRoutes(router, db)
	setupUserAnswerRoutes(router, db, isTest)
	setupQuestionRoutes(router, db, isTest)
	setupQuestionAnswerRoutes(router, db, isTest)
	setupQueryRoutes(router, db, agent, isTest)
	setupTestTopicDataRoutes(router, db, isTest)

	// Setup payment routes
	SetupPaymentRoutes(router, paymentController) // Call the new function to setup payment routes
}
