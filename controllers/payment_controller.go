package controllers

import (
	"log"
	"net/http"
	"os" // Added for Stripe key access, though ideally it comes from config

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"gorm.io/gorm"
)

// PaymentController holds the Stripe client and handles payment-related requests.
type PaymentController struct {
	// StripeClient *stripe.API // It's better to initialize this in main and pass it,
	// but for now, we'll initialize stripe.Key directly in methods or New function
}

// NewPaymentController creates a new PaymentController.
// In a real app, you might pass a pre-configured Stripe client here.
func NewPaymentController(db *gorm.DB) *PaymentController {
	// Initialize Stripe API key from environment variable for now.
	// Ideally, this should come from a config struct passed down from main.go
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Println("Warning: STRIPE_SECRET_KEY environment variable not set.")
	}
	return &PaymentController{}
}

// CreatePaymentIntentRequest defines the expected request body for creating a payment intent.
type CreatePaymentIntentRequest struct {
	Amount   int64  `json:"amount" binding:"required"`   // Amount in the smallest currency unit (e.g., cents)
	Currency string `json:"currency" binding:"required"` // Currency code (e.g., "usd")
}

// CreatePaymentIntent handles the creation of a Stripe PaymentIntent.
func (pc *PaymentController) CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validate amount (e.g., must be greater than a minimum amount)
	if req.Amount <= 0 { // Example: Stripe might have a minimum, e.g. 50 cents for USD
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Validate currency (optional, depends on your app's needs)
	if req.Currency == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Currency is required"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Ensure Stripe key is set. This is a fallback if NewPaymentController wasn't used or key wasn't set there.
	// Best practice is to ensure stripe.Key is set globally once at startup.
	if stripe.Key == "" {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		if stripe.Key == "" {
			log.Println("Critical: STRIPE_SECRET_KEY is not set. Cannot create payment intent.")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment system configuration error."})
			return
		}
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(req.Currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"userID": userID.(string), // Add userID to metadata
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Error creating payment intent: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client_secret": pi.ClientSecret})
}
