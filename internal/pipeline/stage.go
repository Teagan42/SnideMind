package pipeline

type PipelineStage[IN any, OUT any] struct {
	Name   string                     // Name of the stage, used for identification and logging
	Input  <-chan IN                  // Input channel for the stage, where data is received
	Output chan<- *OUT                // Output channel for the stage, where processed data is sent
	Error  chan<- error               // Error channel for the stage, where errors are sent
	Next   *[]PipelineStage[OUT, any] // Pointer to the next stage in the pipeline, allowing chaining of stages
}

type PipelineStageInterface[IN any, OUT any] interface {
	StageName() string
	OnInput(input IN)
}

func (s *PipelineStage[IN, OUT]) StageName() string {
	return s.Name
}

func (s *PipelineStage[IN, OUT]) OnInput(input IN) {
	panic("unimplemented")
}
