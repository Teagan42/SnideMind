package llm

import "go.uber.org/fx"

var Module = fx.Module(
	"LLM Interactation Stage",
	fx.Provide(
		NewLLM,
	),
)
