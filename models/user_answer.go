package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAnswer struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	UserID      string         `json:"user_id" gorm:"not null;index"`     // Foreign key to User model
	QuestionID  string         `json:"question_id" gorm:"not null;index"` // Foreign key to Question model
	AnsweredAt  time.Time      `json:"answered_at" gorm:"not null"`       // Timestamp of when the question was answered
	TimeTaken   int            `json:"time_taken" gorm:"not null"`        // Time taken to answer (in seconds)
	SubjectArea string         `json:"subject_area" gorm:"not null"`      // Subject area of the question
	Attempts    int            `json:"attempts" gorm:"not null"`
}
