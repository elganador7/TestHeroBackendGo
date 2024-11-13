package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	QuestionText string `json:"question_text"`
	Answer       string `json:"answer"`
	Difficulty   int    `json:"difficulty"`
}
