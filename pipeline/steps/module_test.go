package steps

import (
	"testing"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestStepsModule_IsNotNil(t *testing.T) {
	if stepsModule := fx.Option(Module); stepsModule == nil {
		t.Error("steps.Module should not be nil")
	}
}

type TestParams struct {
	fx.In
	StepMap map[string]models.PipelineStepFactory `name:"pipelineStepFactoryMap"`
}

func TestStepsModule_ProvidesPipelineStepFactoryMap(t *testing.T) {

	app := fxtest.New(
		t,
		fx.Provide(func() *config.Config {
			return &config.Config{}
		}),
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		Module,
		fx.Invoke(func(p TestParams) {
			if p.StepMap == nil {
				t.Error("Expected pipelineStepFactoryMap to be provided, got nil")
			}
			if len(p.StepMap) == 0 {
				t.Error("Expected pipelineStepFactoryMap to be provided, got empty map")
			}
			if _, ok := p.StepMap["storeMemory"]; !ok {
				t.Error("Expected 'storeMemory' to be in pipelineStepFactoryMap, got nil")
			}
			if _, ok := p.StepMap["extractTags"]; !ok {
				t.Error("Expected 'extractTags' to be in pipelineStepFactoryMap, got nil")
			}
			if _, ok := p.StepMap["fork"]; !ok {
				t.Error("Expected 'fork' to be in pipelineStepFactoryMap, got nil")
			}
			if _, ok := p.StepMap["llm"]; !ok {
				t.Error("Expected 'llm' to be in pipelineStepFactoryMap, got nil")
			}
			if _, ok := p.StepMap["reduceTools"]; !ok {
				t.Error("Expected 'reduceTools' to be in pipelineStepFactoryMap, got nil")
			}
			if _, ok := p.StepMap["retrieveMemory"]; !ok {
				t.Error("Expected 'retrieveMemory' to be in pipelineStepFactoryMap, got nil")
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
