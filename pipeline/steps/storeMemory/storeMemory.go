package storememory

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
)

type StoreMemory struct {
}

type Params struct {
	fx.In
	Config config.PipelineStepConfig
}

type Result struct {
	fx.Out
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type StoreMemoryFactory struct{}

func (f StoreMemoryFactory) Name() string {
	return "storeMemory"
}
func (f StoreMemoryFactory) Build(config config.PipelineStepConfig) (models.PipelineStep, error) {
	return &StoreMemory{}, nil
}

func NewStoreMemory() (Result, error) {
	return Result{
		Factory: StoreMemoryFactory{},
	}, nil
}

func (s StoreMemory) Name() string {
	return "StoreMemory"
}

func (s StoreMemory) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	// Placeholder for memory storage logic
	// In a real implementation, this would process the input string and store it in memory
	// For now, we'll just return nil to indicate success
	return input, nil
}
