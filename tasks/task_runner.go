package tasks

import (
	"log"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func RunTasks(db *gorm.DB) {
	CreateMaterializedView(db)

	// Initialize the cron scheduler
	c := cron.New()

	// Schedule the task: Every day at 1:00 AM
	_, err := c.AddFunc("0 1 * * *", func() {
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
