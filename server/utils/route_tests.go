package utils

import (
	"net/http"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// mockRoute implements the Route interface for testing.
type mockRouteRoute struct{}

func (m *mockRouteRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
func (m *mockRouteRoute) Methods() []string                                { return []string{"GET"} }
func (m *mockRouteRoute) Pattern() string                                  { return "/mock" }

// mockRouteConstructor returns a Route.
func mockRouteRouteConstructor() Route {
	return &mockRouteRoute{}
}

type testParams struct {
	fx.In
	Routes []Route `group:"testRoutes"`
}

func TestProvideRoutes_ProvidesAnnotatedRoutes(t *testing.T) {
	app := fxtest.New(
		t,
		ProvideRoutes("test", mockRouteRouteConstructor()),
		fx.Invoke(func(p testParams) {
			if len(p.Routes) != 1 {
				t.Fatalf("expected 1 route, got %d", len(p.Routes))
			}
			if p.Routes[0].Pattern() != "/mock" {
				t.Errorf("expected pattern '/mock', got '%s'", p.Routes[0].Pattern())
			}
			if p.Routes[0].Methods()[0] != "GET" {
				t.Errorf("expected method 'GET', got '%s'", p.Routes[0].Methods()[0])
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}

func TestProvideRoutes_EmptyConstructors(t *testing.T) {
	app := fxtest.New(
		t,
		ProvideRoutes("empty"),
		fx.Invoke(func(lc fx.Lifecycle) {
			// No routes should be provided, but app should still start/stop cleanly.
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
