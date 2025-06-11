package llm

import "github.com/teagan42/snidemind/pipeline"

type LLM struct {
}

func NewLLM() *LLM {
	return &LLM{}
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
