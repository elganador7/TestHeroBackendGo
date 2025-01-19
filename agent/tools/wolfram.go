package tools

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/openai/openai-go"
)

var WolframTool = openai.ChatCompletionToolParam{
	Type: openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.F(openai.FunctionDefinitionParam{
		Name:        openai.String("query_wolfram"),
		Description: openai.String("Query the Wolfram Alpha API LLM to compute or verify mathematical expressions"),
		Parameters: openai.F(openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"expression": map[string]string{
					"type":        "string",
					"description": "The mathematical expression to compute or verify and instructions on what to do with it, returns a string detailing the results of the calculation. You should format your request in a way that wolfram can handle and use the result to validate your calculations",
				},
			},
			"required": []string{"expression"},
		}),
	}),
}

type WolframToolObj struct {
	WolframAppID string
}

// queryWolframAlpha sends a mathematical expression to Wolfram Alpha and returns the result
func (w *WolframToolObj) QueryWolframAlpha(expression string) (string, error) {
	appID := w.WolframAppID

	// Encode the query
	query := url.Values{}
	query.Set("input", expression)
	query.Set("output", "json")
	query.Set("appid", appID)
	query.Set("maxchars", "100")

	// Build the API URL
	apiURL := fmt.Sprintf("https://www.wolframalpha.com/api/v1/llm-api?%s", query.Encode())

	// Perform the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to call Wolfram Alpha API: %w", err)
	}
	defer resp.Body.Close()

	responseText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(responseText), nil
}

// Handle tool calls if present
// toolCalls := response.Choices[0].Message.ToolCalls

// counter := 0

// for len(toolCalls) > 0 {
// 	params.Messages.Value = append(params.Messages.Value, response.Choices[0].Message)
// 	for _, toolCall := range toolCalls {
// 		if toolCall.Function.Name == "evaluate_math_expression" {
// 			// Extract the expression from the tool call arguments
// 			var args map[string]interface{}
// 			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
// 				return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("failed to parse tool call arguments: %w", err)
// 			}
// 			expression := args["expression"].(string)

// 			// Query Wolfram Alpha
// 			// wolframResult, err := a.queryWolframAlpha(expression)
// 			// if err != nil {
// 			// 	return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("wolfram Alpha query failed: %w", err)
// 			// }

// 			result, err := tools.EvaluateExpression(expression)
// 			if err != nil {
// 				return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("wolfram Alpha query failed: %w", err)
// 			}

// 			// Return the result to the tool
// 			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, result))
// 			log.Printf("Wolfram response %s", result)
// 		}
// 	}

// 	response, err = a.client.Chat.Completions.New(ctx, params)
// 	if err != nil {
// 		return models.AnswerGeneratorOutputSchema{}, fmt.Errorf("open AI API call failed: %w", err)
// 	}

// 	toolCalls = response.Choices[0].Message.ToolCalls

// 	counter++

// 	log.Printf("Number of tool calls: %d", counter)
// }
