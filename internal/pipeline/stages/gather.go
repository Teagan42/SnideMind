package stages

import (
	"github.com/teagan42/snidemind/internal/pipeline"
)

type GatherStage[IN any, OUT any] struct {
	pipeline.PipelineStage[IN, OUT]
}

func NewGatherStageOptions[IN any, OUT any]() *GatherStage[IN, OUT] {
	return &GatherStage[IN, OUT]{
		PipelineStage: pipeline.PipelineStage[IN, OUT]{
			Name:      "GatherStage",
			Arguments: make(<-chan IN),
			Output:    make(chan *OUT),
			Error:     make(chan error),
		},
	}
}
