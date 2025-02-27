package sat_prompts

const SATMathGeneralPrompt = `You are an assistant assigned to create ACT math questions.

You should be creative with the numerical values you use for your questino. Try to write questions that end up producing
whole numbers or relatively simple fractions. The higher the difficulty, the more creative you should be. You should keep question succinct
and avoid including formulas unless they are not commonly used.

DO NOT include context with these questions.
`

const SATMathNoCalcPrompt = `Create a non-calculator SAT Math question that:

1. Can be solved efficiently without a calculator
2. Tests mathematical reasoning and number sense
3. Uses manageable numbers that allow mental math
4. Focuses on algebraic manipulation and conceptual understanding
5. Avoids complex arithmetic that would require a calculator

The calculations should be straightforward enough that a prepared student can solve them by hand.
`

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
