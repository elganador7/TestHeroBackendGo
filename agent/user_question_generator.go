package agent

import (
	"TestHeroBackendGo/agent/prompts"
	"TestHeroBackendGo/models"
	"fmt"
	"log"
	"math/rand"

	"github.com/google/uuid"
)

const (
	DEFAULT_CORRECT_SCORE          = 1
	DEFAULT_QUESTION_DIFFICULTY    = 0.5
	MINIMUM_QUESTION_DIFFICULTY    = 0.05
	QUESTION_DISTRIBUTION_MODIFIER = 0.1 // Reducing this makes the algorithm give users more questions on weak topics
	// It is added to 1 and then the correct rate is subtracted, meaning 0.05 should yield twice as many questions for
	// a topic in which a student has answered 85% correctly as it does for a topic in which a student has answered 95% correctly
)

// GenerateQuestion generates a question tailored to the current user's performance
func (a *Agent) GenerateRelevantQuestion(questionGenerationInput models.QuestionGeneratorTopicInput) (models.Question, error) {
	// Check if testType and subject are provided
	if questionGenerationInput.Subject == "" || questionGenerationInput.TestType == "" || questionGenerationInput.UserID == "" {
		return models.Question{}, fmt.Errorf("testType, Subject, and UserID are required")
	}

	// Query for all topics, subtopics, and specific topics under the testType and subject
	var testTopics []models.TestTopicData
	if err := a.DB.Where("test_type = ? AND subject = ?", questionGenerationInput.TestType, questionGenerationInput.Subject).Find(&testTopics).Error; err != nil {
		return models.Question{}, fmt.Errorf("failed to fetch test topics, %v", err)
	}

	// Query the user performance summary
	var userPerformance []models.UserPerformanceSummary
	if err := a.DB.Where("user_id = ? AND test_type = ? AND subject = ?", questionGenerationInput.UserID, questionGenerationInput.TestType, questionGenerationInput.Subject).Find(&userPerformance).Error; err != nil {
		return models.Question{}, fmt.Errorf("failed to fetch user performance, %v", err)
	}

	// Create a map of specific topics to performance
	performanceMap := make(map[string]float64)
	questionDifficultyMap := make(map[string]float64)
	for _, performance := range userPerformance {
		if performance.TotalPointsPossible == 0 {
			performanceMap[performance.SpecificTopic] = DEFAULT_CORRECT_SCORE
			questionDifficultyMap[performance.SpecificTopic] = DEFAULT_QUESTION_DIFFICULTY
		}
		performanceMap[performance.SpecificTopic] = (performance.TotalPoints * performance.TotalPoints) / performance.TotalPointsPossible
		questionDifficultyMap[performance.SpecificTopic] = performance.TotalPoints / performance.TotalPointsPossible

	}

	// Create a list of topics with their weights
	var selectedTopic models.TestTopicData
	minWeight := 1000000000000000.0 // Very big number
	difficulty := 0.5

	for _, topic := range testTopics {
		// Get the correct rate for the specific topic
		score, ok := performanceMap[topic.SpecificTopic]
		if !ok {
			// If no data available for this topic, treat it as 50% correct rate
			score = DEFAULT_CORRECT_SCORE
			performanceMap[topic.SpecificTopic] = score
			questionDifficultyMap[topic.SpecificTopic] = DEFAULT_QUESTION_DIFFICULTY
		}

		// Calculate weight: lower correct rate means higher weight
		weight := score * rand.Float64()
		if weight < minWeight {
			minWeight = weight
			if questionDifficultyMap[topic.SpecificTopic] > MINIMUM_QUESTION_DIFFICULTY {
				difficulty = questionDifficultyMap[topic.SpecificTopic]
			} else {
				difficulty = MINIMUM_QUESTION_DIFFICULTY
			}
			selectedTopic = topic
		}
	}

	log.Printf(("Performance Map: %v"), performanceMap)

	log.Printf("Selected topic: %v, Weight: %v, Performance: %v", selectedTopic, minWeight, difficulty)

	// Query for a question based on the selected topic
	question, err := a.GenerateNewQuestionWithTopicData(selectedTopic, difficulty)
	if err != nil {
		return models.Question{}, err
	}

	// Return the selected question
	return question, nil
}

func (a *Agent) GenerateNewQuestionWithTopicData(testTopicData models.TestTopicData, difficulty float64) (models.Question, error) {
	previousQuestionModels := []models.Question{}
	previousQuestionTexts := []string{}

	if err := a.DB.Where("test_topic_id = ?", testTopicData.ID).Find(&previousQuestionModels).Limit(5).Error; err != nil {
		return models.Question{}, err
	}

	for _, question := range previousQuestionModels {
		previousQuestionTexts = append(previousQuestionTexts, question.QuestionText)
	}

	inputSchema := models.NewQuestionGeneratorInputSchema{
		Topic:             testTopicData.Topic,
		Subtopic:          testTopicData.Subtopic,
		SpecificTopic:     testTopicData.SpecificTopic,
		Difficulty:        difficulty,
		PreviousQuestions: previousQuestionTexts,
	}

	specificPrompt, ok := prompts.SubjectTopicPromptMap[testTopicData.TestType][testTopicData.Subject]
	if !ok {
		return models.Question{}, fmt.Errorf("no prompt found for test type %s and subject %s", testTopicData.TestType, testTopicData.Subject)
	}

	systemPrompt := specificPrompt + prompts.BasePromptStructure

	// Call the agent with the system prompt
	questionResponse, err := a.GenerateNewQuestion(inputSchema, systemPrompt)
	if err != nil {
		log.Fatalf("Error generating question: %v", err)
		return models.Question{}, err
	}

	answerResponse, err := a.GenerateAnswer(questionResponse)
	if err != nil {
		log.Fatalf("Error generating answer: %v", err)
		return models.Question{}, err
	}

	optionInput := models.OptionGeneratorInputSchema{
		QuestionText:  questionResponse.QuestionText,
		Explanation:   answerResponse.Explanation,
		CorrectAnswer: answerResponse.CorrectAnswer,
	}

	optionsResponse, err := a.GenerateQuestionOptions(optionInput)
	if err != nil {
		log.Fatalf("Error generating options: %v", err)
		return models.Question{}, err
	}

	// Save the question to the database
	question := models.Question{
		ID:            uuid.NewString(),
		Paragraph:     questionResponse.QuestionContext,
		QuestionText:  questionResponse.QuestionText,
		Options:       optionsResponse.Options,
		EstimatedTime: 60,
		Difficulty:    inputSchema.Difficulty,
		TestTopicID:   testTopicData.ID,
		TestTopic:     testTopicData,
	}

	formattedQuestion, err := a.ValidateMathJaxFormatting(question)
	if err != nil {
		return models.Question{}, err
	}

	question = formattedQuestion.TranslateQuestionOutputSchemaToQuestion(question)

	if err := a.DB.Create(&question).Error; err != nil {
		return models.Question{}, err
	}

	answer := models.QuestionAnswer{
		ID:            uuid.NewString(),
		QuestionID:    question.ID,
		CorrectAnswer: optionsResponse.CorrectAnswer,
		Explanation:   answerResponse.Explanation,
	}

	if err := a.DB.Create(&answer).Error; err != nil {
		return models.Question{}, err
	}

	return question, nil
}
