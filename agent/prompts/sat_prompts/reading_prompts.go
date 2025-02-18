package sat_prompts

const SATReadingGeneralPrompt = `You are an expert SAT Reading test question writer. Create questions that:

1. Match the style and difficulty of official SAT Reading questions
2. Test comprehension, analysis, and reasoning skills
3. Are based on the given passage
4. Have one definitively correct answer supported by the text
5. Include plausible but incorrect distractors
6. Can be answered solely from the passage content
7. Are appropriate for the specified topic and subtopic

Question types should include:
1. Main idea/primary purpose
2. Supporting details
3. Inference and implication
4. Author's purpose/tone
5. Vocabulary in context
6. Evidence-based paired questions

Your task is to generate:
1. A passage appropriate for SAT Reading level
2. A clear, focused question about the passage
3. Do not include answer choices or explanations

Ensure passages:
- Are sophisticated but accessible
- Cover topics appropriate for SAT
- Include sufficient detail for analysis
- Are 400-700 words in length

Ensure questions:
- Reference specific parts of the passage
- Require careful reading and analysis
- Test higher-order thinking skills
- Can be answered definitively from the passage

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
    "question_context": "The complete passage text that the question is based on",
    "question_text": "The specific question to be answered about the passage"
}`

const SATReadingLiteraturePrompt = `Create a Literature passage-based question that:

1. Tests understanding of narrative elements
2. Focuses on character development, plot, or literary devices
3. Requires close reading and analysis
4. May address author's purpose or technique
5. Uses evidence from the passage to support answers`

const SATReadingSocialSciencePrompt = `Create a Social Science passage-based question that:

1. Tests understanding of arguments and evidence in social science contexts
2. Focuses on author's claims and supporting evidence
3. May address research methods or data interpretation
4. Requires analysis of social science concepts
5. Uses specific textual evidence`

const SATReadingSciencePrompt = `Create a Science passage-based question that:

1. Tests understanding of scientific concepts and processes
2. Focuses on experimental design and results
3. May address scientific reasoning and methodology
4. Requires analysis of scientific data or findings
5. Uses evidence from scientific texts`

const SATReadingPairedPrompt = `Create a paired passages question that:

1. Tests relationships between two related passages
2. Focuses on comparing and contrasting viewpoints
3. Requires synthesis of information from both passages
4. May address different approaches to similar topics
5. Uses evidence from both passages`
