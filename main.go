package main

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/config"
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"
	"log"

	"github.com/gin-contrib/cors"
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

	for _, model := range models.AllModels {
		database.DB.AutoMigrate(model)
	}

	router := gin.Default()

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	log.Printf("API key: %s", cfg.OAIAPIKey)

	agent := agent.NewAgent(cfg.OAIAPIKey)

	// parser.ParseJsonData(database.DB)

	routes.SetupRoutes(router, database.DB, agent)

	router.Run(":8080")
}
