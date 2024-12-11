package routes

import (
	"TestHeroBackendGo/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupAuthRoutes(router *gin.Engine, db *gorm.DB) {
	authCtrl := auth.NewAuthController(db)
	authApi := router.Group("/api/auth")

	// Auth routes
	authApi.POST("/register", auth.Register)
	authApi.POST("/login", auth.Login)
	authApi.POST("/google", authCtrl.HandleGoogleAuth)
}
