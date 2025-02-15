package models

type SimilarQuestionGeneratorInputSchema struct {
	Paragraph    string `json:"paragraph" jsonschema_description:"An example paragraph that the original quesiton was based on, you should write your own"`
	QuestionText string `json:"question_text" jsonschema_description:"The existing question text to base the new question on"`
}

type NewQuestionGeneratorInputSchema struct {
	Topic             string   `json:"topic" jsonschema_description:"The topic for the new question, such as Algebra, Chemistry, or Literature"`
	Subtopic          string   `json:"subtopic" jsonschema_description:"The subtopic for the new question, such as Advanced Algebra, Plane Geometry, or Trigonometry"`
	SpecificTopic     string   `json:"specific_topic" jsonschema_description:"The specific topic for the new question, such as Trigonometry, Trigonometric Identities, or Trigonometric Equations"`
	Description       string   `json:"description" jsonschema_description:"The description of the specific_topic for the new question"`
	Difficulty        float64  `json:"difficulty" jsonschema_description:"The difficulty for the new question, a decimal from 0.0 to 1.0 with 1 being the most difficult possible"`
	PreviousQuestions []string `json:"previous_questions" jsonschema_description:"A list of previous questions on this topic"`
}

// OutputSchema defines the structured output for a new question and its answer.
type QuestionGeneratorOutputSchema struct {
	QuestionText string `json:"question_text" jsonschema_description:"The new question text"`
}

type AnswerGeneratorOutputSchema struct {
	Explanation   string `json:"explanation" jsonschema_description:"The explanation for the correct answer"`
	CorrectAnswer string `json:"correct_answer" jsonschema_description:"The correct answer to the new question, given the value of the correct answer"`
}

type OptionGeneratorInputSchema struct {
	QuestionText  string `json:"question_text" jsonschema_description:"The existing question text to base the new question on"`
	Explanation   string `json:"explanation" jsonschema_description:"The explanation for the correct answer"`
	CorrectAnswer string `json:"correct_answer" jsonschema_description:"The correct answer to the new question, given as the index of the correct answer"`
}

type QuestionGeneratorTopicInput struct {
	TestType string `json:"test_type"`
	Subject  string `json:"subject"`
	UserID   string `json:"user_id"`
}

// RawOptionGeneratorOutput represents the array-based output format from the LLM
type OptionGeneratorOutputSchema struct {
	Options       []string `json:"options" jsonschema_description:"The options for the new question as an array"`
	CorrectAnswer string   `json:"correct_answer" jsonschema_description:"The value of the correct answer. Should be a member of the options array"`
}
