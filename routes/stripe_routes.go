package routes

import (
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/utils"

	"github.com/gin-gonic/gin"
)

func SetupStripeRoutes(router *gin.Engine, isTest bool) {
	stripeCtrl := controllers.NewStripeController()

	stripeApi := router.Group("/api/stripe")
	stripeApi.Use(utils.GenerateHandlers(isTest)...)
	{
		stripeApi.GET("/config", stripeCtrl.Config)
		stripeApi.POST("/create-checkout-session", stripeCtrl.CreateCheckoutSession)
		stripeApi.GET("/checkout-session", stripeCtrl.CheckoutSession)
		stripeApi.POST("/customer-portal", stripeCtrl.CustomerPortal)
		stripeApi.POST("/webhook", stripeCtrl.Webhook)
	}
}
