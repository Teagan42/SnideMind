package models

import "github.com/teagan42/snidemind/config"

type PipelineMessage struct {
	Request   *ChatCompletionRequest
	Tags      *[]string               // Tags associated with the message
	Tools     *[]string               // Tools associated with the message
	Prompts   *[]string               // Prompts associated with the message
	Memories  *[]string               // Memories associated with the message
	Knowledge *[]string               // Knowledge associated with the message
	Response  *ChatCompletionResponse // Content of the message
}

type PipelineStep interface {
	Process(previous *[]PipelineStep, in *PipelineMessage) (*PipelineMessage, error) // Process the input data and return the output or an error
	Name() string                                                                    // Name of the stage, used for identification and logging
}

type NewPipelineStep func(config config.PipelineStepConfig) (PipelineStep, error) // Function type for creating a new pipeline step

type PipelineStepFactory interface {
	Name() string                                                 // Name of the factory, used for identification and logging
	Build(config config.PipelineStepConfig) (PipelineStep, error) // Build a new pipeline step from the factory
}
