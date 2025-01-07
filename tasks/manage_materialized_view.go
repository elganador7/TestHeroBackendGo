package tasks

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateMaterializedView(db *gorm.DB) error {
	// Check if the materialized view already exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM pg_matviews WHERE matviewname = 'user_performance_summary'
	);`

	// Execute the query to check existence
	if err := db.Raw(query).Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking for materialized view: %v", err)

	}

	// If the materialized view doesn't exist, create it
	if !exists {
		createQuery := `
			CREATE MATERIALIZED VIEW user_performance_summary AS
			SELECT 
				ua.user_id,
				tt.test_type,
				tt.subject,
				tt.topic,
				tt.subtopic,
				tt.specific_topic,  -- Include specific_topic in the view
				AVG(CASE WHEN ua.attempts > 1 THEN 0 ELSE 1.0 END) AS correct_rate,
				SUM(CASE WHEN ua.attempts > 1 THEN 0 ELSE ua.difficulty*100.0 END) AS total_points,
				SUM(ua.difficulty*100.0) AS total_points_possible,
				COUNT(*) AS question_count
			FROM 
				user_answers ua
			JOIN
				test_topic_data tt ON ua.test_topic_id = tt.id
			GROUP BY 
				ua.user_id, tt.test_type, tt.subject, tt.topic, tt.subtopic, tt.specific_topic;

		`

		// Run the query to create the materialized view
		if err := db.Exec(createQuery).Error; err != nil {
			return fmt.Errorf("error creating materialized view: %v", err)
		}
		fmt.Println("Materialized view 'user_performance_summary' created successfully.")
	} else {
		fmt.Println("Materialized view 'user_performance_summary' already exists.")
	}

	return nil
}

func RefreshMaterializedView(db *gorm.DB) {
	db.Exec("REFRESH MATERIALIZED VIEW user_performance_summary")
}
