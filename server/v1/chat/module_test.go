package chat

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"github.com/teagan42/snidemind/pipeline"
	"github.com/teagan42/snidemind/server/utils"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestChatModule_IsNotNil(t *testing.T) {
	// This test ensures that the chat.Module variable is not nil and is an fx.Option.
	if chatModule := fx.Option(Module); chatModule == nil {
		t.Error("chat.Module should not be nil")
	}
}

type ChatModuleTestParams struct {
	fx.In
	Routes []interface{} `group:"chatRoutes"`
}

func TestChatModule_ProvidesChatRoutes(t *testing.T) {
	var routes []utils.Route
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		fx.Provide(func() *config.Config {
			return &config.Config{}
		}),
		fx.Provide(func() *pipeline.Pipeline {
			return &pipeline.Pipeline{
				Steps:  []models.PipelineStep{},
				Logger: zap.NewNop(),
			}
		}),
		fx.Provide(
			fx.Annotate(
				func() *mux.Router {
					return mux.NewRouter()
				},
				fx.ResultTags(`name:"v1Router"`),
			),
		),
		Module,
		fx.Populate(
			fx.Annotate(
				&routes,
				fx.ParamTags(`group:"chatRoutes"`),
			),
		),
	)
	app.RequireStart()
	if len(routes) == 0 {
		t.Error("Expected chat routes to be provided, but got none")
	}
	if routes[0] == nil {
		t.Error("Expected chat routes to be non-nil, but got nil")
	}
	if route, ok := routes[0].(utils.Route); !ok {
		t.Error("Expected chat routes to be of type utils.Route, but got a different type")
	} else {
		if route.Pattern() != "completions" {
			t.Errorf("Expected route pattern to be '/v1/chat', but got '%s'", route.Pattern())
		}
	}
	app.RequireStop()
}
