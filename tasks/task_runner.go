package tasks

import (
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"TestHeroBackendGo/agent"
	"TestHeroBackendGo/models"
	"TestHeroBackendGo/tasks/topic_data_processor"
)

func RunTasks(db *gorm.DB, agent *agent.Agent, userIdQuestionGenerationChannel chan models.QuestionGeneratorTopicInput) {
	err := CreateMaterializedView(db)
	if err != nil {
		log.Fatalf("Error creating materialized view: %v", err)
	}

	// Path to the root directory where JSON files are stored
	rootDir := "./tasks/topic_data_processor/topic_data"

	// Process the directory and load data into the database
	err = topic_data_processor.ProcessDirectory(rootDir, db)
	if err != nil {
		log.Fatalf("Error processing files: %v", err)
	}

	log.Println("All files processed successfully.")

	// Initialize the cron scheduler
	c := cron.New()

	// Schedule the task: Every day at 1:00 AM
	_, err = c.AddFunc("*/1 * * * *", func() {
		log.Println("Running scheduled task to archive old answers...")
		ArchiveOldAnswersTask(db)
	})
	if err != nil {
		log.Fatalf("Failed to schedule archive task: %v", err)
	}

	// Schedule a task to aggregate user performance every 5 minutes
	_, err = c.AddFunc("*/1 * * * *", func() {
		log.Println("Running scheduled task to aggregate user performance...")
		RefreshMaterializedView(db)
	})
	if err != nil {
		log.Fatalf("Failed to schedule aggregate task: %v", err)
	}

	// Start the scheduler
	c.Start()
}
