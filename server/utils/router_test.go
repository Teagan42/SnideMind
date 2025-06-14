package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// mockRoute implements the Route interface for testing
type mockRoute struct {
	pattern string
	methods []string
}

func (m *mockRoute) Pattern() string {
	return m.pattern
}
func (m *mockRoute) Methods() []string {
	return m.methods
}
func (m *mockRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("mock"))
}

type RouterTestParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Router    *mux.Router `name:"Router"`
}

func TestApiRouteModule_ProvidesRouterAndRegistersRoutes(t *testing.T) {
	mockRoutes := []any{
		func() *mockRoute {
			return &mockRoute{pattern: "test", methods: []string{"GET"}}
		},
	}

	app := fxtest.New(
		t,
		fx.Provide(
			fx.Annotate(
				func() *mux.Router {
					return mux.NewRouter()
				},
				fx.ResultTags(`name:"Router"`),
			),
		),
		ApiRouteModule(ApiRouteModuleParams{
			ModuleName:   "testmod",
			ParentRouter: "",
			Prefix:       "api",
			Routes:       &mockRoutes,
			SubModules:   nil,
		}),
		fx.Invoke(
			func(p RouterTestParams) {
				// Ensure the subrouter is mounted at the correct prefix
				req := httptest.NewRequest("GET", "/api/test", nil)
				rr := httptest.NewRecorder()
				p.Router.ServeHTTP(rr, req)
				require.Equal(t, http.StatusOK, rr.Code)
				require.Equal(t, "mock", rr.Body.String())
			},
		),
	)
	app.RequireStart()
	app.RequireStop()
}

type EmptryRouterTestParams struct {
	fx.In
	fx.Lifecycle
	Router *mux.Router `name:"emptyRouter"`
}

func TestApiRouteModule_WithNoRoutes(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(
			fx.Annotate(
				func() *mux.Router {
					return mux.NewRouter()
				},
				fx.ResultTags(`name:"v1Router"`),
			),
		),
		ApiRouteModule(ApiRouteModuleParams{
			ModuleName:   "empty",
			ParentRouter: "v1",
			Prefix:       "empty",
			Routes:       nil,
			SubModules:   nil,
		}),
		fx.Invoke(
			func(p EmptryRouterTestParams) {
				// Should not panic or register any routes
				req := httptest.NewRequest("GET", "/empty/none", nil)
				rr := httptest.NewRecorder()
				p.Router.ServeHTTP(rr, req)
				require.Equal(t, http.StatusNotFound, rr.Code)
			},
		),
	)
	app.RequireStart()
	app.RequireStop()
}

type SubModuleTestParams struct {
	fx.In
	fx.Lifecycle
	Router *mux.Router `name:"submodRouter"`
}

func TestApiRouteModule_WithSubModules(t *testing.T) {
	subRoute := &mockRoute{pattern: "sub", methods: []string{"GET"}}
	subRoutes := []any{
		func() *mockRoute {
			return subRoute
		},
	}
	subModule := ApiRouteModule(ApiRouteModuleParams{
		ModuleName:   "submod",
		ParentRouter: "testmod",
		Prefix:       "sub",
		Routes:       &subRoutes,
		SubModules:   nil,
	})

	app := fxtest.New(
		t,
		fx.Provide(
			fx.Annotate(
				func() *mux.Router {
					return mux.NewRouter()
				},
				fx.ResultTags(`name:"Router"`),
			),
		),
		ApiRouteModule(ApiRouteModuleParams{
			ModuleName:   "testmod",
			ParentRouter: "",
			Prefix:       "api",
			Routes:       nil,
			SubModules:   &[]fx.Option{subModule},
		}),
		fx.Invoke(
			func(p SubModuleTestParams) {
				req := httptest.NewRequest("GET", "/api/sub/sub", nil)
				rr := httptest.NewRecorder()
				p.Router.ServeHTTP(rr, req)
				require.Equal(t, http.StatusOK, rr.Code)
				require.Equal(t, "mock", rr.Body.String())
			},
		),
	)
	app.RequireStart()
	app.RequireStop()
}
