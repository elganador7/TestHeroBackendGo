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

const (
	FormattingInstructions = `
        Format your response using:
        1. User proper Markdown formatting to organize your response and produce clear, clean output.
        2. If you need to include mathematical expressions, use the following format:
        - Surround all mathematical expressions that are part of a sentence with non-mathematical text with dollar signs like this: $...$
        - Surround all formulas that are on their own line mathematical expressions with double dollar signs like this: $$...$$
        - Ensure all mathematical symbols are properly formatted using Mathjax Compatible LaTeX
    `

	BasePromptStructure = `You are an expert question writer for standardized tests. You should assume you specialize in the test type given.
    
    You will receive in put json object with data about a question you should generate.
        The input will be in the following JSON format:
        {
            "topic": "The main topic area",
            "subtopic": "The specific subtopic",
            "specific_topic": "The specific concept being tested",
            "description": "A description of the specific topic to generate the question about",
            "difficulty": 0.7,  // number between 0 and 1
            "previous_questions": ["..."] // a list of previous questions on this topic, you should avoid questions that are similar 
            to these so that the user gets unique questions
        }
        
        You should treat difficulties between 0 and 0.15 extremely easy, 0.15 to 0.3 easy, 0.3 to 0.4 moderately easy, 0.4 to 0.6 medium difficulty,
        0.6 to 0.7 as moderately hard, 0.7 to 0.9 as hard, and 0.9 to 1 extremely difficult, as difficult as possible.

        Your response should be in the following JSON format:
        {
            "question_context": "Any necessary context for the question (leave empty if not needed)",
            "question_text": "The specific question to be answered"
        }
        
        You should only generate the question and question context (if needed) and nothing else. 

        If the question requires math, you should be creative with the numerical values you use for your questino. For all questions with numerical answers,
        write questions that are answerable with whole numbers or relatively simple fractions.
        
        The higher the difficulty, the more creative you should be. You should keep question succinct
		and avoid including formulas unless they are not commonly used.

        ` + FormattingInstructions
)
