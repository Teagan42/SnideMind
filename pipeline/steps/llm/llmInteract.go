package llm

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/pipeline"
	"go.uber.org/fx"
)

type LLM struct {
	config.LLMConfig
}

type Params struct {
	fx.In
	Config config.PipelineStepConfig
}

type Result struct {
	fx.Out
	Stage *LLM `name:"stage"`
}

func NewLLM(p Params) (*Result, error) {
	return &Result{
		Stage: &LLM{
			LLMConfig: *p.Config.LLM,
		},
	}, nil
}

func (s *LLM) Name() string {
	return "LLM"
}

func (s *LLM) Process(input *pipeline.PipelineMessage) (*pipeline.PipelineMessage, error) {
	// Placeholder for LLM interaction logic
	// In a real implementation, this would process the input string and interact with an LLM
	// For now, we'll just return the input string as output and no error
	return input, nil
}
