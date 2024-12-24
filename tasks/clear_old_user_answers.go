package tasks

import (
	"log"

	"gorm.io/gorm"
)

// ArchiveOldAnswersTask retains only the last 10 answers per user and test_topic_id, moving older ones to an archive
func ArchiveOldAnswersTask(db *gorm.DB) {
	// Step 1: Find user/test_topic_id pairs with more than 10 answers
	var results []struct {
		UserID      string
		TestTopicID string
		Count       int
	}

	err := db.Raw(`
		SELECT user_id, test_topic_id, COUNT(*) as count
		FROM user_answers
		GROUP BY user_id, test_topic_id
		HAVING COUNT(*) > 10
	`).Scan(&results).Error
	if err != nil {
		log.Printf("Error fetching users with excess answers: %v", err)
		return
	}

	// Step 2: For each user/test_topic_id pair, move oldest answers to archive
	for _, result := range results {
		log.Printf("Archiving answers for user_id: %s, test_topic_id: %s (Total Answers: %d)", result.UserID, result.TestTopicID, result.Count)

		// Fetch the IDs of the oldest answers beyond the 10 most recent
		var idsToArchive []string
		if result.Count > 10 {
			err := db.Raw(`
			SELECT id
			FROM user_answers
			WHERE user_id = ? AND test_topic_id = ?
			ORDER BY created_at ASC
			LIMIT ?
		`, result.UserID, result.TestTopicID, result.Count-10).Scan(&idsToArchive).Error
			if err != nil {
				log.Printf("Error fetching IDs to archive: %v", err)
				continue
			}
		} else {
			continue
		}

		if len(idsToArchive) == 0 {
			continue
		}

		// Step 3: Insert the old records into the archive table
		tx := db.Begin()
		if err := tx.Exec(`
			INSERT INTO user_answer_archives (id, user_id, question_id, test_topic_id, time_taken, attempts, difficulty, created_at, updated_at)
			SELECT id, user_id, question_id, test_topic_id, time_taken, attempts, difficulty, created_at, updated_at
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
		log.Printf("Successfully archived %d answers for user_id: %s, test_topic_id: %s", len(idsToArchive), result.UserID, result.TestTopicID)
	}
}
