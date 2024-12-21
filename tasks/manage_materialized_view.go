package tasks

import "gorm.io/gorm"

func CreateMaterializedView(db *gorm.DB) {
	query := `
		CREATE MATERIALIZED VIEW user_performance_summary AS
		SELECT 
			user_id,
			test_type,
			subject,
			topic,
			subtopic,
			AVG(CASE WHEN attempts > 1 THEN 0 ELSE 1.0 END) AS correct_rate
		FROM 
			user_answers
		GROUP BY 
			user_id, test_type, subject, topic, subtopic;
	`

	db.Exec(query)
}

func RefreshMaterializedView(db *gorm.DB) {
	db.Exec("REFRESH MATERIALIZED VIEW user_performance_summary")
}
