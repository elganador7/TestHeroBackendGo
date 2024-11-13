package main

import (
	"TestHeroBackendGo/config"
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
	}

	cfg := config.LoadConfig()

	database.ConnectDatabase(cfg)
	database.DB.AutoMigrate(&models.User{})

	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080")
}
