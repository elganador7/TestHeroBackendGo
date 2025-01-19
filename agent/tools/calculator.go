package tools

import (
	"fmt"
	"log"

	"github.com/Pramod-Devireddy/go-exprtk"
	"github.com/openai/openai-go"
)

// Define Wolfram Alpha tool
var CalculatorTool = openai.ChatCompletionToolParam{
	Type: openai.F(openai.ChatCompletionToolTypeFunction),
	Function: openai.F(openai.FunctionDefinitionParam{
		Name:        openai.String("evaluate_math_expression"),
		Description: openai.String("Evaluate mathematical expressions using the ExprTk library."),
		Parameters: openai.F(openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"expression": map[string]string{
					"type":        "string",
					"description": "A mathematical expression to evaluate, supporting arithmetic, trigonometric functions (e.g., sin, cos), and more. The expression must be valid and formatted correctly.",
				},
			},
			"required": []string{"expression"},
		}),
	}),
}

// Function to evaluate the mathematical expression using go-exprtk
func EvaluateExpression(expression string) (string, error) {
	exprtkObj := exprtk.NewExprtk()
	defer exprtkObj.Delete()

	// Set the expression
	exprtkObj.SetExpression(expression)

	// Compile the expression
	err := exprtkObj.CompileExpression()
	if err != nil {
		log.Default().Printf("Error compiling expression: %s", expression)
		return "", fmt.Errorf("failed to compile expression: %w", err)
	}

	// Evaluate the expression
	result := exprtkObj.GetEvaluatedValue()

	// Return the result as a string
	return fmt.Sprintf("Result: %f", result), nil
}
