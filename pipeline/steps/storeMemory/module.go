package storememory

import "go.uber.org/fx"

var Module = fx.Module(
	"storememory",
	fx.Provide(
		NewStoreMemory,
	),
)
