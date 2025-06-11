package retrievememory

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"Retrieve Memory Stage",
	fx.Provide(
		NewRetrieveMemory,
	),
)
