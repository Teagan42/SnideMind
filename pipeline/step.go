package pipeline

type PipelineStep interface {
	Process(previous *[]PipelineStep, in PipelineMessage) (PipelineMessage, error) // Process the input data and return the output or an error
	Name() string                                                                  // Name of the stage, used for identification and logging
}
