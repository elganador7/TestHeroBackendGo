package utils

import (
	"TestHeroBackendGo/auth"

	"github.com/gin-gonic/gin"
)

func GenerateHandlers(isTest bool) []gin.HandlerFunc {
	if !isTest {
		return []gin.HandlerFunc{auth.JWTAuthMiddleware()}
	}
	return []gin.HandlerFunc{}
}
