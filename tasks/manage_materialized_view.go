package tasks

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateMaterializedView(db *gorm.DB) {
	// Check if the materialized view already exists
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM pg_matviews WHERE matviewname = 'user_performance_summary'
	);`

	// Execute the query to check existence
	if err := db.Raw(query).Scan(&exists).Error; err != nil {
		fmt.Printf("Error checking materialized view existence: %v\n", err)
		return
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
				AVG(CASE WHEN ua.attempts > 1 THEN 0 ELSE 1.0 END) AS correct_rate
			FROM 
				user_answers ua
			JOIN 
				test_topic_data tt ON ua.test_topic_id = tt.id -- Join with TestTopicData on TestTopicID
			GROUP BY 
				ua.user_id, tt.test_type, tt.subject, tt.topic, tt.subtopic, tt.specific_topic; -- Group by specific_topic as well

		`

		// Run the query to create the materialized view
		if err := db.Exec(createQuery).Error; err != nil {
			fmt.Printf("Error creating materialized view: %v\n", err)
			return
		}
		fmt.Println("Materialized view 'user_performance_summary' created successfully.")
	} else {
		fmt.Println("Materialized view 'user_performance_summary' already exists.")
	}
}

func RefreshMaterializedView(db *gorm.DB) {
	db.Exec("REFRESH MATERIALIZED VIEW user_performance_summary")
}
