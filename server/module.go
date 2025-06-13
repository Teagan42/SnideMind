package server

import (
	v1 "github.com/teagan42/snidemind/server/v1"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"Server",
	fx.Provide(
		NewServer,
	),
	v1.Module,
)
