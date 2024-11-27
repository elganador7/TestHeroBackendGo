package models

import "gorm.io/datatypes"

type QuestionGeneratorInputSchema struct {
	Paragraph    string `json:"paragraph" jsonschema_description:"An example paragraph that the original quesiton was based on, you should write your own"`
	QuestionText string `json:"question_text" jsonschema_description:"The existing question text to base the new question on"`
}

// OutputSchema defines the structured output for a new question and its answer.
type QuestionGeneratorOutputSchema struct {
	QuestionText  string `json:"question_text" jsonschema_description:"The new question text"`
	CorrectAnswer string `json:"correct_answer" jsonschema_description:"The correct answer to the new question"`
	Explanation   string `json:"explanation" jsonschema_description:"The explanation for the correct answer"`
}

type OptionGeneratorOutputSchema struct {
	Options       datatypes.JSONMap `json:"options" jsonschema_description:"The options for the new question"`
	CorrectAnswer string            `json:"correct_answer" jsonschema_description:"The correct answer to the new question, given as the correct multiple choice option"`
}
