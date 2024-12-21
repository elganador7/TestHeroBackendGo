package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAnswer struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	UserID      string         `json:"user_id" gorm:"not null;index"`            // Foreign key to User model
	QuestionID  string         `json:"question_id" gorm:"not null;index"`        // Foreign key to Question model
	TimeTaken   int            `json:"time_taken" gorm:"not null"`               // Time taken to answer (in seconds)
	TestTopicID string         `json:"test_topic_id" gorm:"type:uuid;not null"`  // Foreign key reference
	TestTopic   TestTopicData  `json:"test_topic" gorm:"foreignKey:TestTopicID"` // Relation to TestTopicData
	Attempts    int            `json:"attempts" gorm:"not null"`
	Difficulty  float64        `json:"difficulty" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
