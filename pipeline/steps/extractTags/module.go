package extracttags

import "go.uber.org/fx"

var Module = fx.Module(
	"Extract Tags Stage",
	fx.Provide(
		NewExtractTags,
	),
)
