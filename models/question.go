package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Question struct {
	ID            string            `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
	QuestionText  string            `json:"question_text"`
	TestType      string            `json:"test_type"`
	Subject       string            `json:"subject"`
	Topic         string            `json:"topic"`
	Difficulty    float64           `json:"difficulty"`
	Options       datatypes.JSONMap `json:"options" gorm:"type:jsonb"`
	EstimatedTime int               `json:"estimated_time"`
}
