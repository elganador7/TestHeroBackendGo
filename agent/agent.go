package agent

import (
	"TestHeroBackendGo/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Agent struct {
	client *openai.Client
}

// NewAgent initializes and returns a new Agent.
func NewAgent(apiKey string) *Agent {
	client := openai.NewClient(
		option.WithAPIKey(apiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return &Agent{
		client: client,
	}
}

// GenerateSchema generates a JSON schema for a given type.
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

var questionGeneratorOutputSchema = GenerateSchema[models.QuestionGeneratorOutputSchema]()
var optionGeneratorOutputSchema = GenerateSchema[models.OptionGeneratorOutputSchema]()

// GenerateQuestionWithAnswer generates a new question and answer based on an existing question text.
func (a *Agent) GenerateQuestion(input models.QuestionGeneratorInputSchema) (models.QuestionGeneratorOutputSchema, error) {
	ctx := context.Background()

	systemPrompt := `You are an assistant for creating standardized test questions. Expect a JSON input with the following structure:

	{
		"question_text": "..."
	}

	Generate a question that is similar to the one you are provided, you can modify the concept slightly as long as you test a similar topic.
	If the queestion is a math problem, generate a question MathJax rendering in React. Since $ ... $ can conflict with Markdown or certain text processors, \( ... \) is often safer for inline math.
	Please use \( ... \) formatting for inline math even if the question you are provided uses other formatting. 
	
	Take the information from above and output json like the following:

	{
		"question_text": "...",
	}

	Do not respond with anything other than JSON.
	`

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("generate_question_with_answer"),
		Description: openai.F("Generate a new question and its answer based on the input question text"),
		Schema:      openai.F(questionGeneratorOutputSchema),
		Strict:      openai.Bool(true),
	}

	// Prepare the input question
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return models.QuestionGeneratorOutputSchema{}, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Query the Chat Completions API
	response, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(string(inputJSON)),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		return models.QuestionGeneratorOutputSchema{}, fmt.Errorf("API call failed: %w", err)
	}

	// Parse the response into the OutputSchema struct
	var result models.QuestionGeneratorOutputSchema
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
	if err != nil {
		return models.QuestionGeneratorOutputSchema{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

func (a *Agent) GenerateAnswer(input models.QuestionGeneratorOutputSchema) (models.AnswerGeneratorOutputSchema, error) {
	ctx := context.Background()

	systemPrompt := `You are an assistant for creating standardized test questions. Expect a JSON input with the following structure:

	{
		"question_text": "..."
	}

	Think through the question step by step until you reach an answer. You should record the explanation and the final correct answer 
	in the JSON format given below.

	{
		"explanation": "...",
		"correct_answer": "..."
	}

	If the queestion is a math problem, ensure your response and explanation use proper formatting for MathJax rendering in React. 
	Since $ ... $ can conflict with Markdown or certain text processors, \( ... \) is often safer for inline math.
	Please use \( ... \) formatting for inline math even if the question you are provided uses other formatting. 
	
	
	Do not respond with anything other than JSON. You should write the explanation first to ensure that your answer is correct
	and matches your explanation. Do not second guess your answer.
	`

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("generate_question_with_answer"),
		Description: openai.F("Generate a new question and its answer based on the input question text"),
		Schema:      openai.F(questionGeneratorOutputSchema),
		Strict:      openai.Bool(true),
	}

	// Prepare the input question
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Query the Chat Completions API
	response, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(string(inputJSON)),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("API call failed: %w", err)
	}

	// Parse the response into the OutputSchema struct
	var result models.AnswerGeneratorOutputSchema
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
	if err != nil {
		return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// GenerateQuestionWithAnswer generates a new question and answer based on an existing question text.
func (a *Agent) GenerateQuestionOptions(input models.OptionGeneratorInputSchema) (models.OptionGeneratorOutputSchema, error) {
	ctx := context.Background()

	systemPrompt := `You are an assistant for creating standardized test questions. Expect a JSON input with the following structure:

	{
		"question_text": "..."
		"explanation": "..."
		"correct_answer": "..."
	}

	Generate options that make sense in context, including the correct answer as one of the options. Then, output the 
	options as a JSON object and the correct answer as the letter option that it is. Each of the options should be compatible with
	MathJax rendering in React. Since $ ... $ can conflict with Markdown or certain text processors, \( ... \) is often safer for inline math.
	Make sure to wrap all math with that formatting, even if it is not used in the question.

	{
		"options": {
			"A": "...",
			"B": "...",
			"C": "...",
			"D": "..."
		},
		"correct_option": "A",
	}

	Do not respond with anything other than JSON. You should write the explanation first to ensure that your answer is correct
	and matches your explanation.
	`

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("generate_question_with_answer"),
		Description: openai.F("Generate a new question and its answer based on the input question text"),
		Schema:      openai.F(optionGeneratorOutputSchema),
		Strict:      openai.Bool(true),
	}

	// Prepare the input question
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return models.OptionGeneratorOutputSchema{}, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Query the Chat Completions API
	response, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(string(inputJSON)),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4oMini),
	})
	if err != nil {
		return models.OptionGeneratorOutputSchema{}, fmt.Errorf("API call failed: %w", err)
	}

	// Parse the response into the OutputSchema struct
	var result models.OptionGeneratorOutputSchema
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
	if err != nil {
		return models.OptionGeneratorOutputSchema{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
