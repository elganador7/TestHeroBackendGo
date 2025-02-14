package prompts

import (
	"TestHeroBackendGo/agent/prompts/act_prompts"
	"TestHeroBackendGo/agent/prompts/gmat_prompts"
	"TestHeroBackendGo/agent/prompts/gre_prompts"
	"TestHeroBackendGo/agent/prompts/lsat_prompts"
	"TestHeroBackendGo/agent/prompts/mcat_prompts"
	"TestHeroBackendGo/agent/prompts/sat_prompts"
)

var PromptMap = map[string]map[string]string{
	"ACT": {
		"Math": act_prompts.ActMathPrompt,
	},
	"SAT": {
		"Math":       sat_prompts.SATMathGeneralPrompt,
		"NoCalc":     sat_prompts.SATMathNoCalcPrompt,
		"Calc":       sat_prompts.SATMathCalcPrompt,
		"GridIn":     sat_prompts.SATMathGridInPrompt,
		"Reading":    sat_prompts.SATReadingGeneralPrompt,
		"Literature": sat_prompts.SATReadingLiteraturePrompt,
		"Science":    sat_prompts.SATReadingSciencePrompt,
		"Social":     sat_prompts.SATReadingSocialSciencePrompt,
		"Paired":     sat_prompts.SATReadingPairedPrompt,
	},
	"LSAT": {
		"Logical":    lsat_prompts.LSATLogicalPrompt,
		"Reading":    lsat_prompts.LSATReadingPrompt,
		"Analytical": lsat_prompts.LSATAnalyticalPrompt,
	},
	"GRE": {
		"Quantitative": gre_prompts.GREQuantitativePrompt,
		"Verbal":       gre_prompts.GREVerbalPrompt,
		"Analytical":   gre_prompts.GREAnalyticalPrompt,
	},
	"GMAT": {
		"Quantitative": gmat_prompts.GMATQuantitativePrompt,
		"Verbal":       gmat_prompts.GMATVerbalPrompt,
		"Integrated":   gmat_prompts.GMATIntegratedPrompt,
	},
	"MCAT": {
		"Chemical":      mcat_prompts.MCATChemicalPhysicalPrompt,
		"CARS":          mcat_prompts.MCATCarsPrompt,
		"Biological":    mcat_prompts.MCATBiologicalPrompt,
		"Psychological": mcat_prompts.MCATPsychologicalPrompt,
	},
	"Default_Prompt": {},
}
