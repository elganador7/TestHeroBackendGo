package routes

import (
	"TestHeroBackendGo/controllers"
	"TestHeroBackendGo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupTestTopicDataRoutes(router *gin.Engine, db *gorm.DB, isTest bool) {
	testTopicCtrl := controllers.TestTopicDataController{
		DB: db,
	}

	testTopicApi := router.Group("/api/test-topic-data")

	// Protected routes (Require authentication)
	testTopicApi.Use(utils.GenerateHandlers(isTest)...)
	{
		testTopicApi.GET("/", testTopicCtrl.ListTestTopicData)                                 // List all TestTopicData
		testTopicApi.GET("/:id", testTopicCtrl.GetTestTopicData)                               // Get single TestTopicData by ID
		testTopicApi.POST("/create", testTopicCtrl.CreateTestTopicData)                        // Create new TestTopicData
		testTopicApi.PUT("/:id", testTopicCtrl.UpdateTestTopicData)                            // Update TestTopicData by ID
		testTopicApi.DELETE("/:id", testTopicCtrl.DeleteTestTopicData)                         // Delete TestTopicData by ID
		testTopicApi.GET("/test-type/:test_type", testTopicCtrl.ListByTestType)                // List by TestType
		testTopicApi.GET("/subject/:subject", testTopicCtrl.ListBySubject)                     // List by Subject
		testTopicApi.GET("/topic/:topic", testTopicCtrl.ListByTopic)                           // List by Topic
		testTopicApi.GET("/subtopic/:subtopic", testTopicCtrl.ListBySubtopic)                  // List by Subtopic
		testTopicApi.GET("/specific-topic/:specific_topic", testTopicCtrl.ListBySpecificTopic) // List by SpecificTopic
	}
}
