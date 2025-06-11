package storememory

import "go.uber.org/fx"

var Module = fx.Module(
	"Pipeline Stages",
	fx.Provide(
		NewStoreMemory,
	),
)
