package main

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/config"
	"TestHeroBackendGo/database"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/routes"
	"TestHeroBackendGo/tasks"
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

	for _, model := range models.AllModels {
		database.DB.AutoMigrate(model)
	}

	router := gin.Default()

	// Configure CORS middleware
	// router.Use(
	// 	cors.New(
	// 		cors.Config{
	// 			AllowOrigins:     []string{"http://localhost:5173*", "https://app.testscorehero.com*", "https://testherobackendgo-production.up.railway.app/*"},
	// 			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 			AllowCredentials: true,
	// 			MaxAge:           12 * 60 * 60,
	// 			AllowWildcard:    true,
	// 			AllowHeaders: []string{
	// 				"Origin",
	// 				"Content-Type",
	// 				"Authorization",
	// 				"Accept",
	// 				"X-Requested-With",
	// 			},
	// 			ExposeHeaders: []string{"Content-Length", "Content-Type"},
	// 			// AllowAllOrigins: true,
	// 		},
	// 	),
	// )
	router.Use(corsMiddleware())

	userIdGenerationQuestionChannel := make(chan models.QuestionGeneratorTopicInput)

	agent := agent.NewAgent(cfg.OAIAPIKey, database.DB)

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

// CORS middleware function definition
func corsMiddleware() gin.HandlerFunc {
	// Define allowed origins as a comma-separated string
	allowedOrigins := []string{"http://localhost:5173", "https://app.testscorehero.com", "https://testherobackendgo-production.up.railway.app/"}

	// Return the actual middleware handler function
	return func(c *gin.Context) {
		// Function to check if a given origin is allowed
		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		}

		// Get the Origin header from the request
		origin := c.Request.Header.Get("Origin")
		log.Printf("Origin: %s", origin)

		// Check if the origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			log.Println("Origin is allowed")
			// If the origin is allowed, set CORS headers in the response
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		}

		// Handle preflight OPTIONS requests by aborting with status 204
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Call the next handler
		c.Next()
	}
}
