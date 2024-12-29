// tasks/question_tasks.go
package tasks

import (
	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/models"
	"log"

	"gorm.io/gorm"
)

// StartQuestionWorker initializes the worker for processing question generation.
func MonitorTestTopicChannel(db *gorm.DB, agent *agent.Agent, userIdQuestionGenerationChannel chan models.QuestionGeneratorTopicInput) {
	for {
		// Await the next value from the channel
		questionTestTopic := <-userIdQuestionGenerationChannel

		// Process the received value
		log.Printf("Received topic input: %+v", questionTestTopic)

		// Example: You could add logic to handle question generation here
		question, err := agent.GenerateRelevantQuestion(questionTestTopic)
		if err != nil {
			log.Printf("Failed to generate question: %v", err)
		}

		queuedQuestion := models.QueuedQuestion{
			UserID:     questionTestTopic.UserID,
			QuestionID: question.ID,
			Question:   question,
			TestType:   questionTestTopic.TestType,
			Subject:    questionTestTopic.Subject,
		}

		db.Create(&queuedQuestion)

	}
}
