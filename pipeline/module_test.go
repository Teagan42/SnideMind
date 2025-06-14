package pipeline

import (
	"testing"

	"github.com/teagan42/snidemind/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestPipelineModule_IsNotNil(t *testing.T) {
	if pipelineModule := fx.Option(Module); pipelineModule == nil {
		t.Error("pipeline.Module should not be nil")
	}
}

type PipelineTestParams struct {
	fx.In
	Pipeline *Pipeline
}

func TestPipelineModule_ProvidesPipeline(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		fx.Provide(func() *config.Config {
			return &config.Config{
				Server: config.ServerConfig{
					Port: 8080,
					Bind: nil,
				},
				MCPServers: nil,
				Pipeline: &config.PipelineConfig{
					Steps: []config.PipelineStepConfig{
						{
							Type: "storeMemory",
						},
					},
				},
			}
		}),
		Module,
		fx.Invoke(func(p PipelineTestParams) {
			if p.Pipeline == nil {
				t.Error("Expected Pipeline to be provided, got nil")
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
