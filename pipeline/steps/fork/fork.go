package fork

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
)

type ForkPipelineStage struct {
}

type Params struct {
	fx.In
	Config config.PipelineStepConfig
}

type Result struct {
	fx.Out
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type ForkPipelineStageFactory struct{}

func (f ForkPipelineStageFactory) Name() string {
	return "fork"
}
func (f ForkPipelineStageFactory) Build(config config.PipelineStepConfig) (models.PipelineStep, error) {
	return &ForkPipelineStage{}, nil
}

func NewFork() (Result, error) {
	return Result{
		Factory: ForkPipelineStageFactory{},
	}, nil
}

func (f ForkPipelineStage) Name() string {
	return "Fork Stage"
}

func (f ForkPipelineStage) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	// Placeholder for fork logic
	// In a real implementation, this would process the input and fork it into multiple outputs
	// For now, we'll just return the input as is
	return input, nil
}
