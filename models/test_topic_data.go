package models

type TestTopicData struct {
	ID            string `json:"id" gorm:"type:uuid;primaryKey"`
	TestType      string `json:"test_type" gorm:"not null"`
	Subject       string `json:"subject" gorm:"not null"`
	Topic         string `json:"topic" gorm:"not null"`
	Subtopic      string `json:"subtopic" gorm:"not null"`
	SpecificTopic string `json:"specific_topic" gorm:"not null"`
	Description   string `json:"description" gorm:"not null"`
}
