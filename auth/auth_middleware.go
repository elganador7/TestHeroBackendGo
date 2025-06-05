package auth

import (
	"log"
	"net/http"
	"strings"

	"TestHeroBackendGo/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitializeDB sets the database connection for the auth package
func InitializeDB(database *gorm.DB) {
	db = database
}

// IsPayingCustomer checks if a user with the given email is a paying customer
func IsPayingCustomer(email string) (bool, error) {
	var customer models.StripeCustomer
	result := db.Where("email = ? AND delinquent = ?", email, false).First(&customer)
	if result.Error != nil {
		log.Printf("Error checking if user is paying customer: %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	log.Printf("User is paying customer: %v", customer)
	return true, nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["userID"])
		c.Set("email", claims["email"])
		c.Next()
	}
}

// PayingCustomerMiddleware checks if the authenticated user is a paying customer
func PayingCustomerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get email from context (set by JWTAuthMiddleware)
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in token"})
			c.Abort()
			return
		}

		// Check if user is a paying customer
		isPaying, err := IsPayingCustomer(email.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking payment status"})
			c.Abort()
			return
		}

		if !isPaying {
			c.JSON(http.StatusForbidden, gin.H{"error": "This feature requires an active subscription"})
			c.Abort()
			return
		}

		c.Next()
	}
}
