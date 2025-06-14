package storememory

import (
	"testing"

	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestStoreMemoryModule_IsNotNil(t *testing.T) {
	// This test ensures that the storememory.Module variable is not nil and is an fx.Option.
	if storeMemoryModule := fx.Option(Module); storeMemoryModule == nil {
		t.Error("storememory.Module should not be nil")
	}
}

type StoreMemoryTestParams struct {
	fx.In
	Factory []models.PipelineStepFactory `group:"pipelineStepFactory"`
}

func TestStoreMemoryModule_ProvidesStoreMemory(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		Module,
		fx.Invoke(func(p StoreMemoryTestParams) {
			if p.Factory == nil {
				t.Error("Expected StoreMemory to be provided, got nil")
			}
			if len(p.Factory) == 0 {
				t.Error("Expected StoreMemory to be provided, got empty slice")
			}
			if p.Factory[0].Name() != "storeMemory" {
				t.Errorf("Expected factory name to be 'StoreMemory', got '%s'", p.Factory[0].Name())
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
