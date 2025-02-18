package prompts

var (
	MathJaxFormatter = `
		You will receive a JSON object containing a set of fields defiining a question, an answer, 
		and options.

		Ensure that the structured data input is formatted in a way that will render properly in using a 
		markdown renderer in React with MathJax to render embedded LaTeX. 
		
		If there are errors in the formatting, modify the object if needed and return it in the same JSON format that you
		receive it in.
	`
)

const FormattingInstructions = `Format your response using:
1. Proper Markdown formatting for:
   - Lists (use - for bullets)
   - Emphasis (use * or _)
   - Tables (use | and -)
   - Line breaks (use double space)
2. LaTeX for all mathematical expressions:
   - Use $...$ for inline math formulas (e.g., $x + y = 5$)
   - Use $$...$$ for block math formulas (e.g., $$\frac{d}{dx}(x^2) = 2x$$)
   - Ensure all mathematical symbols are properly formatted`

const BasePromptStructure = `Your task is to generate:
1. A context appropriate for the question type (if needed)
2. A clear, focused question about the context
3. Do not include answer choices or explanations

` + FormattingInstructions + `

The input will be in the following JSON format:
{
    "topic": "The main topic area",
    "subtopic": "The specific subtopic",
    "specific_topic": "The specific concept being tested",
    "difficulty": 0.7,  // number between 0 and 1
    "previous_questions": ["..."] // a list of previous questions on this topic
}

Your response should be in the following JSON format:
{
    "question_context": "Any necessary context for the question (leave empty if not needed)",
    "question_text": "The specific question to be answered"
}`
