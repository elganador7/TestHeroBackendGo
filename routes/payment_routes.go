package routes

import (
	"TestHeroBackendGo/auth"        // Adjust import path if necessary
	"TestHeroBackendGo/controllers" // Adjust import path if necessary

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm" // Not strictly needed for this specific route setup if DB isn't directly used by these routes
)

// SetupPaymentRoutes configures the routes for payment processing.
// It takes the Gin engine and the PaymentController as input.
func SetupPaymentRoutes(router *gin.Engine, paymentController *controllers.PaymentController) {
	// Create a new route group for payments, e.g., /api/v1/payments
	// Applying JWTAuthMiddleware to the entire group ensures all payment routes are protected.
	paymentRoutes := router.Group("/payments")
	paymentRoutes.Use(auth.JWTAuthMiddleware()) // Assuming JWTAuthMiddleware is correctly defined in your auth package
	{
		// Route for creating a payment intent
		// POST /payments/create-payment-intent
		paymentRoutes.POST("/create-payment-intent", paymentController.CreatePaymentIntent)

		// You can add more payment-related routes here in the future,
		// for example, handling webhooks from Stripe, retrieving payment history, etc.
	}
}
