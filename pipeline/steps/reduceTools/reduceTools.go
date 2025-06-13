package reducetools

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/schema"
	"go.uber.org/fx"
)

type ReduceTools struct {
}

type Params struct {
	fx.In
	Config config.PipelineStepConfig
}

type Result struct {
	fx.Out
	Factory schema.PipelineStepFactory `group:"pipelineStepFactory"`
}

type ReduceToolsFactory struct{}

func (f ReduceToolsFactory) Name() string {
	return "reduceTools"
}
func (f ReduceToolsFactory) Build(config config.PipelineStepConfig) (schema.PipelineStep, error) {
	return &ReduceTools{}, nil
}

func NewReduceTools() (Result, error) {
	return Result{
		Factory: ReduceToolsFactory{},
	}, nil
}

func (s ReduceTools) Name() string {
	return "ReduceTools"
}

func (s ReduceTools) Process(previous *[]schema.PipelineStep, input *schema.PipelineMessage) (*schema.PipelineMessage, error) {
	// Placeholder for tool reduction logic
	// In a real implementation, this would process the input and reduce tools accordingly
	// For now, we'll just return a zero value of OUT and no error
	return input, nil
}
