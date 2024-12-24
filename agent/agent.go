package agent

import (
	"TestHeroBackendGo/agent/prompts"
	"TestHeroBackendGo/agent/prompts/base_prompts"
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
var answerGeneratorOutputSchema = GenerateSchema[models.AnswerGeneratorOutputSchema]()
var optionGeneratorOutputSchema = GenerateSchema[models.OptionGeneratorOutputSchema]()
var rawQuestionOutputSchema = GenerateSchema[models.QuestionOutputSchema]()

// GenerateQuestionWithAnswer generates a new question and answer based on an existing question text.
func (a *Agent) GenerateSimilarQuestion(input models.SimilarQuestionGeneratorInputSchema, systemPrompt string) (models.QuestionGeneratorOutputSchema, error) {
	ctx := context.Background()

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
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.Float(0.4),
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

// GenerateQuestionWithAnswer generates a new question and answer based on an existing question text.
func (a *Agent) GenerateNewQuestion(input models.NewQuestionGeneratorInputSchema, systemPrompt string) (models.QuestionGeneratorOutputSchema, error) {
	ctx := context.Background()

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("generate_question_with_answer"),
		Description: openai.F("Generate a new question and its answer based on the input json"),
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
		Model:       openai.F(openai.ChatModelGPT4oMini),
		Temperature: openai.Float(0.6),
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

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("generate_question_with_answer"),
		Description: openai.F("Generate a new question and its answer based on the input question text"),
		Schema:      openai.F(answerGeneratorOutputSchema),
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
			openai.SystemMessage(prompts.AnswerGeneratorSystemPrompt),
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

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("generate_options_for_question"),
		Description: openai.F("Generate options for a question based on the question text and answer in json"),
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
			openai.SystemMessage(prompts.OptionGeneratorSystemPrompt),
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

// GenerateQuestionWithAnswer generates a new question and answer based on an existing question text.
func (a *Agent) ValidateMathJaxFormatting(input models.Question) (models.QuestionOutputSchema, error) {
	ctx := context.Background()

	// Create a structured output parameter
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("validate_math_jax_formatting"),
		Description: openai.F("Generate options for a question based on the question text and answer in json"),
		Schema:      openai.F(rawQuestionOutputSchema),
		Strict:      openai.Bool(true),
	}

	// Prepare the input question
	inputJSON, err := json.Marshal(input.TranslateQuestionToQuestionOutputSchema())
	if err != nil {
		return models.QuestionOutputSchema{}, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Query the Chat Completions API
	response, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(base_prompts.MathJaxFormatter),
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
		return models.QuestionOutputSchema{}, fmt.Errorf("API call failed: %w", err)
	}

	// Parse the response into the OutputSchema struct
	var result models.QuestionOutputSchema
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
	if err != nil {
		return models.QuestionOutputSchema{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
