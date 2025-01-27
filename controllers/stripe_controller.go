package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	portalsession "github.com/stripe/stripe-go/v74/billingportal/session"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/webhook"
)

type StripeController struct {
	publishableKey string
	freePriceId    string
	basicPriceId   string
	proPriceId     string
}

func NewStripeController(
	publishableKey string,
	freePriceId string,
	basicPriceId string,
	proPriceId string,

) *StripeController {
	return &StripeController{
		publishableKey: publishableKey,
		freePriceId:    freePriceId,
		basicPriceId:   basicPriceId,
		proPriceId:     proPriceId,
	}
}

func (ctrl *StripeController) Config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"publishableKey": ctrl.publishableKey,
		"basicPrice":     ctrl.basicPriceId,
		"proPrice":       ctrl.proPriceId,
	})
}

func (ctrl *StripeController) CreateCheckoutSession(c *gin.Context) {
	var req struct {
		PriceID string `json:"priceId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(os.Getenv("DOMAIN") + "/success.html?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(os.Getenv("DOMAIN") + "/canceled.html"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(req.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(params)
	if err != nil {
		writeJSON(c, nil, err)
		return
	}
	c.Redirect(http.StatusSeeOther, s.URL)
}

func (ctrl *StripeController) CheckoutSession(c *gin.Context) {
	sessionID := c.Query("sessionId")
	s, err := session.Get(sessionID, nil)
	writeJSON(c, s, err)
}

func (ctrl *StripeController) CustomerPortal(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := session.Get(req.SessionID, nil)
	if err != nil {
		writeJSON(c, nil, err)
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(s.Customer.ID),
		ReturnURL: stripe.String(os.Getenv("DOMAIN")),
	}
	ps, err := portalsession.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer portal session"})
		return
	}
	c.Redirect(http.StatusSeeOther, ps.URL)
}

func (ctrl *StripeController) Webhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to construct event"})
		return
	}

	if event.Type == "checkout.session.completed" {
		log.Println("Checkout session completed")
	}
	c.Status(http.StatusOK)
}

func writeJSON(c *gin.Context, v interface{}, err error) {
	if err != nil {
		msg := err.Error()
		var serr *stripe.Error
		if errors.As(err, &serr) {
			msg = serr.Msg
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, v)
}
