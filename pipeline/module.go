package pipeline

import (
	"github.com/teagan42/snidemind/pipeline/steps"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"Pipeline",
	steps.Module,
	fx.Provide(
		NewPipeline,
	),
)
