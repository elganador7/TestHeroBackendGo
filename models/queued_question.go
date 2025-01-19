package models

type QueuedQuestion struct {
	ID         int      `gorm:"primaryKey"`
	UserID     string   `json:"user_id"`
	QuestionID string   `json:"question_id"`
	Question   Question `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Foreign key to Question table
	TestType   string   `json:"test_type"`
	Subject    string   `json:"subject"`
}
