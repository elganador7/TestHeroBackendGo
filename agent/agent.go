package agent

import (
	"TestHeroBackendGo/agent/prompts"
	"TestHeroBackendGo/agent/prompts/base_prompts"
	"TestHeroBackendGo/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"gorm.io/gorm"
)

type Agent struct {
	client *openai.Client
	DB     *gorm.DB
}

// NewAgent initializes and returns a new Agent.
func NewAgent(apiKey string, db *gorm.DB) *Agent {
	client := openai.NewClient(
		option.WithAPIKey(apiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	return &Agent{
		client: client,
		DB:     db,
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

	// Define Wolfram Alpha tool
	wolframTool := openai.ChatCompletionToolParam{
		Type: openai.F(openai.ChatCompletionToolTypeFunction),
		Function: openai.F(openai.FunctionDefinitionParam{
			Name:        openai.String("query_wolfram"),
			Description: openai.String("Query the Wolfram Alpha API to compute or verify mathematical expressions"),
			Parameters: openai.F(openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"expression": map[string]string{
						"type":        "string",
						"description": "The mathematical expression to compute or verify.",
					},
				},
				"required": []string{"expression"},
			}),
		}),
	}

	// Set up the API parameters
	params := openai.ChatCompletionNewParams{
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
		Tools: openai.F([]openai.ChatCompletionToolParam{wolframTool}),
	}

	// Query the Chat Completions API
	response, err := a.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("open AI API call failed: %w", err)
	}

	// Handle tool calls if present
	toolCalls := response.Choices[0].Message.ToolCalls
	log.Println("printing toolCalls")
	log.Println("toolCalls: ", toolCalls)
	log.Printf("toolCalls: %v", toolCalls)

	if len(toolCalls) > 0 {
		for _, toolCall := range toolCalls {
			if toolCall.Function.Name == "query_wolfram" {
				// Extract the expression from the tool call arguments
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("failed to parse tool call arguments: %w", err)
				}
				expression := args["expression"].(string)

				// Query Wolfram Alpha
				wolframResult, err := queryWolframAlpha(expression)
				if err != nil {
					return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("wolfram Alpha query failed: %w", err)
				}

				// Return the result to the tool
				params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, wolframResult))
			}
		}
	}

	// Parse the response into the OutputSchema struct
	var result models.AnswerGeneratorOutputSchema
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &result)
	if err != nil {
		return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// queryWolframAlpha sends a mathematical expression to Wolfram Alpha and returns the result
func queryWolframAlpha(expression string) (string, error) {
	appID := "your-wolfram-alpha-app-id" // Replace with your Wolfram Alpha App ID

	// Encode the query
	query := url.Values{}
	query.Set("input", expression)
	query.Set("format", "plaintext")
	query.Set("output", "JSON")
	query.Set("appid", appID)

	// Build the API URL
	apiURL := fmt.Sprintf("https://www.wolframalpha.com/api/v1/llm-apiy?%s", query.Encode())

	log.Printf("Wolfram Alpha URL: %s", apiURL)

	// Perform the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to call Wolfram Alpha API: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var wolframResponse struct {
		QueryResult struct {
			Pods []struct {
				Title   string `json:"title"`
				Subpods []struct {
					Plaintext string `json:"plaintext"`
				} `json:"subpods"`
			} `json:"pods"`
		} `json:"queryresult"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wolframResponse); err != nil {
		return "", fmt.Errorf("failed to parse Wolfram Alpha API response: %w", err)
	}

	// Extract the computation result
	for _, pod := range wolframResponse.QueryResult.Pods {
		if pod.Title == "Result" {
			for _, subpod := range pod.Subpods {
				return subpod.Plaintext, nil
			}
		}
	}

	return "", fmt.Errorf("no result found in Wolfram Alpha response")
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
