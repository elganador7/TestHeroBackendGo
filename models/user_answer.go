package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAnswer struct {
	ID         string         `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	UserID     string         `json:"user_id" gorm:"not null;index"`     // Foreign key to User model
	QuestionID string         `json:"question_id" gorm:"not null;index"` // Foreign key to Question model
	TimeTaken  int            `json:"time_taken" gorm:"not null"`        // Time taken to answer (in seconds)
	TestType   string         `json:"test_type"`
	Subject    string         `json:"subject"`
	Topic      string         `json:"topic"`
	Subtopic   string         `json:"subtopic"`
	Attempts   int            `json:"attempts" gorm:"not null"`
}
