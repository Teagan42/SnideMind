package llm

import "go.uber.org/fx"

var Module = fx.Module(
	"llm",
	fx.Provide(
		NewLLM,
	),
)
