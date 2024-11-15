package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAnswer struct {
	gorm.Model              // Adds fields ID, CreatedAt, UpdatedAt, DeletedAt
	UserID        uint      `json:"user_id" gorm:"not null;index"`     // Foreign key to User model
	QuestionID    uint      `json:"question_id" gorm:"not null;index"` // Foreign key to Question model
	AnsweredAt    time.Time `json:"answered_at" gorm:"not null"`       // Timestamp of when the question was answered
	TimeTaken     int       `json:"time_taken" gorm:"not null"`        // Time taken to answer (in seconds)
	IsCorrect     bool      `json:"is_correct" gorm:"not null"`        // Whether the answer was correct
	SubjectArea   string    `json:"subject_area" gorm:"not null"`      // Subject area of the question
	AnswerDetails string    `json:"answer_details" gorm:"type:jsonb"`  // Additional metadata in JSON format
}
