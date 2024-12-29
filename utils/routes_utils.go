package utils

import (
	"TestHeroBackendGo/auth"

	"github.com/gin-gonic/gin"
)

func GenerateHandlers(isTest bool) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{}
	if !isTest {
		handlers = append(handlers, auth.JWTAuthMiddleware())
	}
	return handlers
}
