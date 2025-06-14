package retrievememory

import (
	"testing"

	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestRetrieveMemoryModule_IsNotNil(t *testing.T) {
	if retrieveMemoryModule := fx.Option(Module); retrieveMemoryModule == nil {
		t.Error("retrievememory.Module should not be nil")
	}
}

type RetrieveMemoryTestParams struct {
	fx.In
	Factory []models.PipelineStepFactory `group:"pipelineStepFactory"`
}

func TestRetrieveMemoryModule_ProvidesRetrieveMemory(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		Module,
		fx.Invoke(func(p RetrieveMemoryTestParams) {
			if p.Factory == nil {
				t.Error("Expected RetrieveMemory to be provided, got nil")
			}
			if len(p.Factory) == 0 {
				t.Error("Expected RetrieveMemory to be provided, got empty slice")
			}
			if p.Factory[0].Name() != "retrieveMemory" {
				t.Errorf("Expected factory name to be 'RetrieveMemory', got '%s'", p.Factory[0].Name())
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
