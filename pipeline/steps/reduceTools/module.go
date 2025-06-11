package reducetools

import "go.uber.org/fx"

var Module = fx.Module(
	"Reduce Tools Stage",
	fx.Provide(
		NewReduceTools,
	),
)
