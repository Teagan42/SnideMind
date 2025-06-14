package reducetools

import (
	"testing"

	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestReduceToolsModule_IsNotNil(t *testing.T) {
	if reducetoolsModule := fx.Option(Module); reducetoolsModule == nil {
		t.Error("reducetools.Module should not be nil")
	}
}

type ReduceToolsTestParams struct {
	fx.In
	Factory []models.PipelineStepFactory `group:"pipelineStepFactory"`
}

func TestReduceToolsModule_ProvidesReduceTools(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		Module,
		fx.Invoke(func(p ReduceToolsTestParams) {
			if p.Factory == nil {
				t.Error("Expected ReduceTools to be provided, got nil")
			}
			if len(p.Factory) == 0 {
				t.Error("Expected ReduceTools to be provided, got empty slice")
			}
			if p.Factory[0].Name() != "reduceTools" {
				t.Errorf("Expected factory name to be 'ReduceTools', got '%s'", p.Factory[0].Name())
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
