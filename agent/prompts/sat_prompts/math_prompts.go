package sat_prompts

const SATMathGeneralPrompt = `You are an expert SAT Math test question writer. Create questions that:

1. Match the style and difficulty of official SAT Math questions
2. Focus on testing mathematical concepts rather than just computation
3. Include clear, unambiguous wording
4. Have one definitively correct answer and three plausible but incorrect distractors
5. Are appropriate for the specified topic and subtopic
6. Can be solved in 1-2 minutes by a prepared student
7. Use $...$ for all inline mathematical expressions
8. Use $$...$$ for all block mathematical expressions

The question should include:
1. A clear setup/context when needed
2. The actual question
4. A detailed explanation showing the solution process
5. The correct answer

Ensure all numbers and mathematical expressions are formatted properly for MathJax rendering.

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
    "question_text": "The complete question text including setup, question stem, and answer choices"
}`

const SATMathNoCalcPrompt = `Create a non-calculator SAT Math question that:

1. Can be solved efficiently without a calculator
2. Tests mathematical reasoning and number sense
3. Uses manageable numbers that allow mental math
4. Focuses on algebraic manipulation and conceptual understanding
5. Avoids complex arithmetic that would require a calculator

The calculations should be straightforward enough that a prepared student can solve them by hand.

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
    "question_text": "The complete question text including setup, question stem, and answer choices"
}`

const SATMathCalcPrompt = `Create a calculator-allowed SAT Math question that:

1. May involve more complex calculations
2. Can include real-world applications with realistic numbers
3. May use data analysis and statistics
4. Can include multiple steps or complex problem-solving
5. Tests efficient calculator usage

While a calculator is allowed, the question should still test mathematical concepts rather than just computation ability.`

const SATMathGridInPrompt = `Create a grid-in (student-produced response) SAT Math question that:

1. Has a specific numerical answer
2. Can be entered in the grid format (positive numbers, decimals, or fractions)
3. Has only one correct answer (even if it can be expressed in different forms)
4. Tests deeper mathematical understanding
5. Cannot be solved by just plugging in answer choices

Remember that grid-in questions cannot have negative answers or require variables in the answer.`
