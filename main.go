package main

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/config"
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"
	"TestHeroBackendGo/tasks"
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
	router.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"http://localhost:5173", "https://app.testscorehero.com"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * 60 * 60,
			},
		),
	)

	userIdGenerationQuestionChannel := make(chan models.QuestionGeneratorTopicInput)

	agent := agent.NewAgent(cfg.OAIAPIKey, database.DB, cfg.WolframAppID)

	// Should this be run from the main thread?
	go tasks.MonitorTestTopicChannel(database.DB, agent, userIdGenerationQuestionChannel)

	// Start Tasks
	tasks.RunTasks(database.DB, agent, userIdGenerationQuestionChannel)

	routes.SetupRoutes(router, database.DB, agent, false)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.Run(":8080")
}
