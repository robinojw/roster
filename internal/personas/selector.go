package personas

import (
	"github.com/robinojw/roster/internal/analyser"
)

const (
	frameworkNextJS    = "Next.js"
	frameworkNuxt      = "Nuxt"
	frameworkSvelteKit = "SvelteKit"
	frameworkDjango    = "Django"
	frameworkRails     = "Rails"
)

var corePersonas = map[string]bool{
	"architect": true,
	"reviewer":  true,
	"test":      true,
	"docs":      true,
	"security":  true,
}

var frontendPersonas = map[string]bool{
	"design":        true,
	"accessibility": true,
}

var signalPersonas = map[string]func(*analyser.RepoSignals) bool{
	"devops": func(signals *analyser.RepoSignals) bool {
		return signals.CIProvider != "" || signals.HasDocker
	},
	"api": func(signals *analyser.RepoSignals) bool {
		return hasAnyFramework(signals, "Express", "Fastify", "Gin", "Fiber", "Echo",
			frameworkDjango, "Flask", frameworkRails, frameworkNextJS, frameworkNuxt)
	},
	"data": func(signals *analyser.RepoSignals) bool {
		return hasAnyFramework(signals, frameworkDjango, frameworkRails) ||
			hasAnyLanguage(signals, "Python", "Ruby", "Java", "Kotlin", "C#")
	},
	"performance": func(signals *analyser.RepoSignals) bool {
		return hasFrontend(signals) || signals.HasE2E ||
			hasAnyFramework(signals, frameworkNextJS, frameworkNuxt, frameworkSvelteKit)
	},
}

func Select(all []Persona, signals *analyser.RepoSignals) []Persona {
	included := buildIncludedSet(signals)

	var selected []Persona
	for _, p := range all {
		if included[p.ID] {
			selected = append(selected, p)
		}
	}
	return selected
}

func buildIncludedSet(signals *analyser.RepoSignals) map[string]bool {
	included := make(map[string]bool)

	for id := range corePersonas {
		included[id] = true
	}

	if hasFrontend(signals) {
		for id := range frontendPersonas {
			included[id] = true
		}
	}

	for id, condition := range signalPersonas {
		if condition(signals) {
			included[id] = true
		}
	}

	return included
}

func hasFrontend(signals *analyser.RepoSignals) bool {
	hasDesignSignals := signals.HasDesignSystem || signals.HasStorybook
	if hasDesignSignals {
		return true
	}
	hasFrontendFramework := hasAnyFramework(signals, "React", "Vue", "Svelte", "Angular",
		frameworkNextJS, frameworkNuxt, frameworkSvelteKit)
	hasFrontendLanguage := hasAnyLanguage(signals, "TypeScript", "JavaScript")
	return hasFrontendFramework || hasFrontendLanguage
}

func hasAnyFramework(signals *analyser.RepoSignals, names ...string) bool {
	for _, fw := range signals.Frameworks {
		for _, name := range names {
			if fw == name {
				return true
			}
		}
	}
	return false
}

func hasAnyLanguage(signals *analyser.RepoSignals, names ...string) bool {
	for _, lang := range signals.Languages {
		for _, name := range names {
			if lang == name {
				return true
			}
		}
	}
	return false
}
