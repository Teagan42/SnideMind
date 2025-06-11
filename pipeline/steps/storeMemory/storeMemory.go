package storememory

import "github.com/teagan42/snidemind/pipeline"

type StoreMemory struct {
}

func NewStoreMemory() *StoreMemory {
	return &StoreMemory{}
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
