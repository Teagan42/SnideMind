package retrievememory

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"retrievememory",
	fx.Provide(
		NewRetrieveMemory,
	),
)
