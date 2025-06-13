package reducetools

import "go.uber.org/fx"

var Module = fx.Module(
	"reducetools",
	fx.Provide(
		NewReduceTools,
	),
)
