package tasks

import "gorm.io/gorm"

func CreateMaterializedView(db *gorm.DB) {
	db.Exec(`
		CREATE MATERIALIZED VIEW user_performance_summary AS
		SELECT user_id, subject_area, AVG(is_correct::int) AS correct_rate
		FROM user_answers
		GROUP BY user_id, subject_area;
	`)
}

func RefreshMaterializedView(db *gorm.DB) {
	db.Exec("REFRESH MATERIALIZED VIEW user_performance_summary")
}
