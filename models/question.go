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
	TestTopicID   string            `json:"test_topic_id" gorm:"type:uuid;not null"`  // Foreign key reference
	TestTopic     TestTopicData     `json:"test_topic" gorm:"foreignKey:TestTopicID"` // Relation to TestTopicData
	Difficulty    float64           `json:"difficulty"`
	Options       datatypes.JSONMap `json:"options" gorm:"type:jsonb"`
	EstimatedTime int               `json:"estimated_time"`
	Paragraph     string            `json:"paragraph"`
}
