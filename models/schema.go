package models

import "gorm.io/datatypes"

type SimilarQuestionGeneratorInputSchema struct {
	Paragraph    string `json:"paragraph" jsonschema_description:"An example paragraph that the original quesiton was based on, you should write your own"`
	QuestionText string `json:"question_text" jsonschema_description:"The existing question text to base the new question on"`
}

type NewQuestionGeneratorInputSchema struct {
	TestType   string  `json:"test_type" jsonschema_description:"The test type for the new question, such as the SAT or ACT"`
	Subject    string  `json:"subject" jsonschema_description:"The subject for the new question, such as Math, Science, or English"`
	Topic      string  `json:"topic" jsonschema_description:"The topic for the new question, such as Algebra, Chemistry, or Literature"`
	Subtopic   string  `json:"subtopic" jsonschema_description:"The subtopic for the new question, such as Advanced Algebra, Plane Geometry, or Trigonometry"`
	Difficulty float64 `json:"difficulty" jsonschema_description:"The difficulty for the new question, a decimal from 0.0 to 1.0"`
}

// OutputSchema defines the structured output for a new question and its answer.
type QuestionGeneratorOutputSchema struct {
	QuestionText string `json:"question_text" jsonschema_description:"The new question text"`
}

type AnswerGeneratorOutputSchema struct {
	Explanation   string `json:"explanation" jsonschema_description:"The explanation for the correct answer"`
	CorrectAnswer string `json:"correct_answer" jsonschema_description:"The correct answer to the new question, given as the correct multiple choice option"`
}

type OptionGeneratorInputSchema struct {
	QuestionText  string `json:"question_text" jsonschema_description:"The existing question text to base the new question on"`
	Explanation   string `json:"explanation" jsonschema_description:"The explanation for the correct answer"`
	CorrectAnswer string `json:"correct_answer" jsonschema_description:"The correct answer to the new question, given as the correct multiple choice option"`
}

type OptionGeneratorOutputSchema struct {
	Options       datatypes.JSONMap `json:"options,omitempty" jsonschema_description:"The options for the new question"`
	CorrectOption string            `json:"correct_option" jsonschema_description:"The correct answer to the new question, given as the correct multiple choice option"`
}
