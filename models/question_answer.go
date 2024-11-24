package models

import (
	"time"

	"gorm.io/gorm"
)

type QuestionAnswer struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	QuestionID    string         `json:"question_id" gorm:"not null;index"`
	CorrectAnswer string         `json:"correct_answer"`
	Explanation   string         `json:"explanation"`
}
