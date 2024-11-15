package routes

import (
	"TestHeroBackendGo/auth"

	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(router *gin.Engine) {
	authApi := router.Group("/api/auth")

	// Auth routes
	authApi.POST("/register", auth.Register)
	authApi.POST("/login", auth.Login)
}
