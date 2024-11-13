package main

import (
	"TestHeroBackendGo/config"
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
	}

	cfg := config.LoadConfig()

	// Connect to database
	database.ConnectDatabase(cfg)

	// Set up router
	router := gin.Default()
	routes.SetupRoutes(router)

	// Start server
	router.Run(":8080")
}
