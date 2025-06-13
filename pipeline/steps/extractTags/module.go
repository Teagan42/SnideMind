package extracttags

import "go.uber.org/fx"

var Module = fx.Module(
	"extracttags",
	fx.Provide(
		NewExtractTags,
	),
)
