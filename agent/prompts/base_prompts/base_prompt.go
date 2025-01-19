package base_prompts

var (
	MathJaxFormatter = `
		You will receive a JSON object containing a set of fields defiining a question, an answer, 
		and options.

		Ensure that the structured data input is formatted in a way that will render properly in
		MathJax in React. Modify the object if needed and return it in the same JDON format that you
		receive it in.

		Ensure that you never use $ ... $ for math. Use \( ... \) instead to brackt math in all cases
		because that is required to render properly in MathJax in React.
	`
)
