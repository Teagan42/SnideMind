package models

import (
	"net/http"

	"maps"

	"github.com/teagan42/snidemind/config"
)

type PipelineMessage struct {
	Request        *ChatCompletionRequest
	Tags           *map[string]string      // Tags associated with the message
	Tools          *[]string               // Tools associated with the message
	Prompts        *[]string               // Prompts associated with the message
	Memories       *[]string               // Memories associated with the message
	Knowledge      *[]string               // Knowledge associated with the message
	ResponseWriter http.ResponseWriter     // Content of the message
	Response       *ChatCompletionResponse // Response from the message
}

func (p *PipelineMessage) Combine(message *PipelineMessage) {
	if message.Tags != nil {
		if p.Tags == nil {
			p.Tags = &map[string]string{}
		}
		maps.Copy((*p.Tags), *message.Tags)
	}
	if message.Tools != nil {
		if p.Tools == nil {
			p.Tools = &[]string{}
		}
		*p.Tools = append(*p.Tools, (*message.Tools)...)
	}
	if message.Prompts != nil {
		if p.Prompts == nil {
			p.Prompts = &[]string{}
		}
		*p.Prompts = append(*p.Prompts, (*message.Prompts)...)
	}
	if message.Memories != nil {
		if p.Memories == nil {
			p.Memories = &[]string{}
		}
		*p.Memories = append(*p.Memories, (*message.Memories)...)
	}
	if message.Knowledge != nil {
		if p.Knowledge == nil {
			p.Knowledge = &[]string{}
		}
		*p.Knowledge = append(*p.Knowledge, (*message.Knowledge)...)
	}
}

type PipelineStep interface {
	Process(previous *[]PipelineStep, in *PipelineMessage) (*PipelineMessage, error) // Process the input data and return the output or an error
	Name() string                                                                    // Name of the stage, used for identification and logging
}

type NewPipelineStep func(config config.PipelineStepConfig, stepFactories map[string]PipelineStepFactory) (PipelineStep, error) // Function type for creating a new pipeline step

type PipelineStepFactory interface {
	Name() string                                                                                               // Name of the factory, used for identification and logging
	Build(config config.PipelineStepConfig, stepFactories map[string]PipelineStepFactory) (PipelineStep, error) // Build a new pipeline step from the factory
}
