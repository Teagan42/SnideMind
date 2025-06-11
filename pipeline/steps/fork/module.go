package fork

import "go.uber.org/fx"

var Module = fx.Module(
	"fork",
	fx.Provide(
		NewFork,
	),
)
