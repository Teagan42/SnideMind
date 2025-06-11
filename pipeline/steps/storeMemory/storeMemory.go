package storememory

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/pipeline"
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
	Stage *StoreMemory `name:"stage"`
}

func NewStoreMemory(p Params) (*Result, error) {
	return &Result{
		Stage: &StoreMemory{},
	}, nil
}

func (s *StoreMemory) Name() string {
	return "StoreMemory"
}

func (s *StoreMemory) Process(input *pipeline.PipelineMessage) (*pipeline.PipelineMessage, error) {
	// Placeholder for memory storage logic
	// In a real implementation, this would process the input string and store it in memory
	// For now, we'll just return nil to indicate success
	return input, nil
}
