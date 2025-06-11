package pipeline

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/schema"
	"go.uber.org/fx"
)

type Pipeline struct {
	Steps []PipelineStep // All Steps in the pipeline
}

type PipelineMessage struct {
	Request   schema.ChatCompletionRequest
	Tags      *[]string // Tags associated with the message
	Tools     *[]string // Tools associated with the message
	Prompts   *[]string // Prompts associated with the message
	Memories  *[]string // Memories associated with the message
	Knowledge *[]string // Knowledge associated with the message
}

type Params struct {
	fx.In
	Config config.Config
}

func NewPipeline(p Params) *Pipeline {
	return &Pipeline{
		Steps: make([]PipelineStep, 0),
	}
}

func (p *Pipeline) AddStep(stage PipelineStep, index *int) {
	if index == nil || *index < 0 || *index > len(p.Steps) {
		p.Steps = append(p.Steps, stage)
	} else {
		p.Steps = append(p.Steps[:*index], append([]PipelineStep{stage}, p.Steps[*index:]...)...)
	}
}

func (p *Pipeline) Process(input PipelineMessage) (PipelineMessage, error) {
	var err error
	var previous *[]PipelineStep
	for _, stage := range p.Steps {
		if input, err = stage.Process(previous, input); err != nil {
			return *new(PipelineMessage), err // Return zero value of OUT and the error
		}
		previous = &[]PipelineStep{stage} // Update previous to the current stage
	}

	return input, nil // Return the processed output
}
