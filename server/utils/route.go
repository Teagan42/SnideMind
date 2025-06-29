package utils

import (
	"fmt"
	"net/http"

	"go.uber.org/fx"
)

type Route interface {
	http.Handler
	Methods() []string
	Pattern() string
}

func ProvideRoutes(routeGroup string, routeConstructors ...interface{}) fx.Option {
	var annotatedRoutes []interface{}
	for _, constructor := range routeConstructors {
		annotatedRoutes = append(annotatedRoutes, fx.Annotate(
			constructor,
			fx.As(new(Route)),
			fx.ResultTags(`group:"`+routeGroup+`Routes"`),
		))
	}
	fmt.Printf("Annotated Routes: %v\n", annotatedRoutes)
	return fx.Provide(
		annotatedRoutes...,
	)
}
