package models

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	QuestionText  string         `json:"question_text"`
	TestTopicID   string         `json:"test_topic_id" gorm:"type:uuid;not null"`  // Foreign key reference
	TestTopic     TestTopicData  `json:"test_topic" gorm:"foreignKey:TestTopicID"` // Relation to TestTopicData
	Difficulty    float64        `json:"difficulty"`
	Options       []string       `json:"options" gorm:"serializer:json"`
	EstimatedTime int            `json:"estimated_time"`
	Paragraph     string         `json:"paragraph"`
}

type QuestionOutputSchema struct {
	QuestionText  string   `json:"question_text"`
	Difficulty    float64  `json:"difficulty"`
	Options       []string `json:"options" jsonschema_description:"The options for the new question as an array"`
	EstimatedTime int      `json:"estimated_time"`
	Paragraph     string   `json:"paragraph"`
}

func (question Question) TranslateQuestionToQuestionOutputSchema() QuestionOutputSchema {
	return QuestionOutputSchema{
		QuestionText:  question.QuestionText,
		Difficulty:    question.Difficulty,
		Options:       question.Options,
		EstimatedTime: question.EstimatedTime,
		Paragraph:     question.Paragraph,
	}
}

func (formattedQuestion QuestionOutputSchema) TranslateQuestionOutputSchemaToQuestion(originalQuestion Question) Question {
	return Question{
		ID:            originalQuestion.ID,
		QuestionText:  formattedQuestion.QuestionText,
		TestTopicID:   originalQuestion.TestTopicID,
		TestTopic:     originalQuestion.TestTopic,
		Difficulty:    formattedQuestion.Difficulty,
		Options:       formattedQuestion.Options,
		EstimatedTime: formattedQuestion.EstimatedTime,
		Paragraph:     formattedQuestion.Paragraph,
		CreatedAt:     originalQuestion.CreatedAt,
		UpdatedAt:     originalQuestion.UpdatedAt,
		DeletedAt:     originalQuestion.DeletedAt,
	}
}
