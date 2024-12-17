package tasks

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// UserAnswer represents the main table
type UserAnswer struct {
	ID         int       `gorm:"primaryKey"`
	UserID     int       `gorm:"column:user_id"`
	TestType   string    `gorm:"column:test_type"`
	QuestionID int       `gorm:"column:question_id"`
	IsCorrect  bool      `gorm:"column:is_correct"`
	Attempts   int       `gorm:"column:attempts"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

// ArchiveOldAnswersTask moves the oldest answers to an archive table if there are more than 200 answers for a user/test type
func ArchiveOldAnswersTask(db *gorm.DB) {
	// Step 1: Find user/test_type pairs with more than 200 answers
	var results []struct {
		UserID   int
		TestType string
		Count    int
	}

	err := db.Raw(`
		SELECT user_id, test_type, COUNT(*) as count
		FROM user_answers
		GROUP BY user_id, test_type
		HAVING COUNT(*) > 200
	`).Scan(&results).Error
	if err != nil {
		log.Printf("Error fetching users with excess answers: %v", err)
		return
	}

	// Step 2: For each user/test_type pair, move oldest answers to archive
	for _, result := range results {
		log.Printf("Archiving answers for user_id: %d, test_type: %s (Total Answers: %d)", result.UserID, result.TestType, result.Count)

		// Fetch the IDs of the oldest answers beyond the 200 most recent
		var idsToArchive []int
		err := db.Raw(`
			SELECT id
			FROM user_answers
			WHERE user_id = ? AND test_type = ?
			ORDER BY created_at ASC
			LIMIT ?
		`, result.UserID, result.TestType, result.Count-200).Scan(&idsToArchive).Error
		if err != nil {
			log.Printf("Error fetching IDs to archive: %v", err)
			continue
		}

		if len(idsToArchive) == 0 {
			continue
		}

		// Step 3: Insert the old records into the archive table
		tx := db.Begin()
		if err := tx.Exec(`
			INSERT INTO user_answers_archive (user_id, test_type, question_id, is_correct, attempts, created_at)
			SELECT user_id, test_type, question_id, is_correct, attempts, created_at
			FROM user_answers
			WHERE id IN ?
		`, idsToArchive).Error; err != nil {
			tx.Rollback()
			log.Printf("Error archiving answers: %v", err)
			continue
		}

		// Step 4: Delete the old records from the main table
		if err := tx.Exec(`
			DELETE FROM user_answers
			WHERE id IN ?
		`, idsToArchive).Error; err != nil {
			tx.Rollback()
			log.Printf("Error deleting archived answers: %v", err)
			continue
		}

		// Commit the transaction
		tx.Commit()
		log.Printf("Successfully archived %d answers for user_id: %d, test_type: %s", len(idsToArchive), result.UserID, result.TestType)
	}
}
