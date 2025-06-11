package pipeline

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"Pipeline",
	fx.Provide(
		NewPipeline,
	),
)
