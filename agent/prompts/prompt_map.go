package prompts

import "TestHeroBackendGo/agent/prompts/act_prompts"

var PromptMap = map[string]map[string]string{
	"ACT": {
		"Math": act_prompts.ActMathPrompt,
	},
	"SAT": {
		"Math": act_prompts.ActMathPrompt,
	},
	"Default_Prompt": {},
}
