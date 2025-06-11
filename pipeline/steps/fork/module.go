package fork

import "go.uber.org/fx"

var Module = fx.Module(
	"Fork Pipeline Stages",
	fx.Provide(
		NewFork,
	),
)
