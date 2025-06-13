package utils

import (
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type PathPrefix string

type Router interface {
	Prefix() string
	RegisterRoutes() []Route
}

type ApiRouteModuleParams struct {
	ModuleName   string       `name:"moduleName"`
	ParentRouter string       `name:"parentRouter"`
	Prefix       string       `name:"prefix"`
	Routes       *[]any       `group:"routes"`
	SubModules   *[]fx.Option `group:"submodules"`
}

func (p *ApiRouteModuleParams) GetRoutes() []any {
	if p.Routes == nil {
		return []any{}
	}
	return *p.Routes
}

func ProvidePathPrefix(name string, prefix string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() PathPrefix {
				return PathPrefix(prefix)
			},
			fx.ResultTags(`name:"`+name+`Prefix"`),
		),
	)
}

func RegisterGroupedRoutes(router *mux.Router, prefix PathPrefix, routes []Route) {
	for _, route := range routes {
		pattern := route.Pattern()
		if !strings.HasPrefix(pattern, "/") {
			pattern = "/" + pattern
		}
		router.Handle(pattern, route).Methods(route.Methods()...)
	}
}

func ApiRouteModule(params ApiRouteModuleParams) fx.Option {
	var subModules []fx.Option
	if params.SubModules != nil {
		subModules = *params.SubModules
	}
	return fx.Module(
		params.ModuleName,
		ProvidePathPrefix(params.ModuleName, params.Prefix),
		ProvideRoutes(params.ModuleName, params.GetRoutes()...),
		fx.Provide(
			fx.Annotate(
				func(parent *mux.Router, prefix PathPrefix) *mux.Router {
					sub := parent.PathPrefix("/" + string(prefix)).Subrouter()
					sub.UseEncodedPath()
					return sub
				},
				fx.ParamTags(`name:"`+params.ParentRouter+`Router"`, `name:"`+params.ModuleName+`Prefix"`),
				fx.ResultTags(`name:"`+params.ModuleName+`Router"`),
			),
		),
		fx.Invoke(
			fx.Annotate(
				RegisterGroupedRoutes,
				fx.ParamTags(`name:"`+params.ModuleName+`Router"`, `name:"`+params.ModuleName+`Prefix"`, `group:"`+params.ModuleName+`Routes"`),
			),
		),
		fx.Options(subModules...),
	)
}
