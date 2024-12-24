package models

type UserPerformanceSummary struct {
	UserID              string  `json:"user_id" gorm:"column:user_id"`
	TestType            string  `json:"test_type" gorm:"column:test_type"`
	Subject             string  `json:"subject" gorm:"column:subject"`
	Topic               string  `json:"topic" gorm:"column:topic"`
	Subtopic            string  `json:"subtopic" gorm:"column:subtopic"`
	SpecificTopic       string  `json:"specific_topic" gorm:"column:specific_topic"`
	CorrectRate         float64 `json:"correct_rate" gorm:"column:correct_rate"`
	TotalPoints         float64 `json:"total_points" gorm:"column:total_points"`
	TotalPointsPossible float64 `json:"total_points_possible" gorm:"column:total_points_possible"`
}

// TableName overrides the default table name for GORM
func (UserPerformanceSummary) TableName() string {
	return "user_performance_summary" // Name of your materialized view
}
