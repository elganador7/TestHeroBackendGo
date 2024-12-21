package parser

import (
	"TestHeroBackendGo/models"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// JSONQuestion represents the JSON structure for each question
type JSONQuestion struct {
	ID         string `json:"id"`
	Domain     string `json:"domain"`
	Difficulty string `json:"difficulty"`
	Question   struct {
		Choices struct {
			A string `json:"A"`
			B string `json:"B"`
			C string `json:"C"`
			D string `json:"D"`
		} `json:"choices"`
		QuestionText  string `json:"question"`
		Explanation   string `json:"explanation"`
		CorrectAnswer string `json:"correct_answer"`
		Paragraph     string `json:"paragraph"`
	} `json:"question"`
}

func ParseJsonData(db *gorm.DB) {
	// Open the JSON file
	file, err := os.Open("./questions.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Parse the JSON file
	var data struct {
		Math    []JSONQuestion `json:"math"`
		English []JSONQuestion `json:"english"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Insert questions into the database
	for _, q := range data.Math {
		difficulty, err := parseDifficulty(q.Difficulty)
		if err != nil {
			log.Printf("Skipping question %s due to difficulty parse error: %v", q.ID, err)
			continue
		}

		question := models.Question{
			ID:           uuid.New().String(),
			QuestionText: q.Question.QuestionText,
			Difficulty:   difficulty,
			Options: datatypes.JSONMap{
				"A": q.Question.Choices.A,
				"B": q.Question.Choices.B,
				"C": q.Question.Choices.C,
				"D": q.Question.Choices.D,
			},
			EstimatedTime: 60, // Example estimated time
			Paragraph:     q.Question.Paragraph,
		}

		if err := db.Create(&question).Error; err != nil {
			log.Printf("Failed to insert question %s: %v", q.ID, err)
		}

		answer := models.QuestionAnswer{
			ID:            uuid.New().String(),
			QuestionID:    question.ID,
			CorrectAnswer: q.Question.CorrectAnswer,
			Explanation:   q.Question.Explanation,
		}

		if err := db.Create(&answer).Error; err != nil {
			log.Printf("Failed to insert answer for question %s: %v", q.ID, err)
		}
	}

	fmt.Println("Questions added successfully!")
}

// parseDifficulty converts difficulty from string to float64
func parseDifficulty(difficulty string) (float64, error) {
	switch difficulty {
	case "Easy":
		return 0.3, nil
	case "Medium":
		return 0.6, nil
	case "Hard":
		return 0.9, nil
	default:
		return 0, fmt.Errorf("invalid difficulty level: %s", difficulty)
	}
}
