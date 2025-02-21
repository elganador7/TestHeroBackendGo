package topic_data_processor

import (
	"TestHeroBackendGo/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpecificTopic struct {
	SpecificTopic string `json:"specific_topic"`
	Description   string `json:"description"`
}

type Subtopic struct {
	Subtopic       string          `json:"subtopic"`
	SpecificTopics []SpecificTopic `json:"specific_topics"`
	Description    string          `json:"description"`
}

type Topic struct {
	Topic     string     `json:"topic"`
	Subtopics []Subtopic `json:"subtopics"`
}

type Subject struct {
	Subject string  `json:"subject"`
	Topics  []Topic `json:"topics"`
}

type TestData struct {
	TestType string    `json:"test_type"`
	Subjects []Subject `json:"subjects"`
}

// Check if a TestTopicData entry already exists
func checkIfExists(db *gorm.DB, testType, subject, topic, subtopic, specificTopic string) bool {
	var existing models.TestTopicData
	result := db.Where("test_type = ? AND subject = ? AND topic = ? AND subtopic = ? AND specific_topic = ?",
		testType, subject, topic, subtopic, specificTopic).First(&existing)
	return result.Error == nil
}

// Process a single JSON file and insert data into the database
func processFile(filePath string, db *gorm.DB) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	// Parse the JSON data into Go structure
	var testData TestData
	if err := json.Unmarshal(data, &testData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %v", filePath, err)
	}

	// Insert data into the database
	for _, subject := range testData.Subjects {
		for _, topic := range subject.Topics {
			for _, subtopic := range topic.Subtopics {
				for _, specificTopic := range subtopic.SpecificTopics {
					if exists := checkIfExists(db, testData.TestType, subject.Subject, topic.Topic, subtopic.Subtopic, specificTopic.SpecificTopic); exists {
						continue
					}

					testTopic := models.TestTopicData{
						ID:            uuid.New().String(),
						TestType:      testData.TestType,
						Subject:       subject.Subject,
						Topic:         topic.Topic,
						Subtopic:      subtopic.Subtopic,
						SpecificTopic: specificTopic.SpecificTopic,
						Description:   specificTopic.Description,
					}

					// Insert each row into the database
					if err := db.Create(&testTopic).Error; err != nil {
						return fmt.Errorf("failed to insert data for file %s: %v", filePath, err)
					}
				}
			}
		}
	}

	return nil
}

// Process all JSON files in the provided directory and load data into the database
func ProcessDirectory(rootDir string, db *gorm.DB) error {
	// Walk through the directory structure and process JSON files
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a JSON file
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Process each JSON file
			fmt.Printf("Processing file: %s\n", path)
			if err := processFile(path, db); err != nil {
				return fmt.Errorf("error processing file %s: %v", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path %v: %v", rootDir, err)
	}

	return nil
}
