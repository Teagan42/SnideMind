package stages

import "github.com/teagan42/snidemind/internal/pipeline"

type ProcessorStage[IN any, OUT any] struct {
	pipeline.PipelineStage[IN, OUT]
	Process func(IN) (OUT, error) // Function to process input data
}

type NewProcessorStageOptions[IN any, OUT any] struct {
	Name             string                // Name of the stage, used for identification and logging
	Process          func(IN) (OUT, error) // Function to process input data
	ArgumentsChannel <-chan IN             // Optional input channel for the stage
	OutputChannel    *chan *OUT            // Optional output channel for the stage
	ErrorChannel     *chan error           // Optional error channel for the stage
}

func NewProcessorStage[IN any, OUT any](options NewProcessorStageOptions[IN, OUT]) *ProcessorStage[IN, OUT] {
	stage := ProcessorStage[IN, OUT]{
		PipelineStage: pipeline.PipelineStage[IN, OUT]{
			Name:      options.Name,
			Arguments: options.ArgumentsChannel,
		},
		Process: options.Process,
	}
	if options.OutputChannel != nil {
		stage.Output = *options.OutputChannel
	} else {
		stage.Output = make(chan *OUT) // Create a default channel if none provided
	}
	if options.ErrorChannel != nil {
		stage.Error = *options.ErrorChannel
	} else {
		stage.Error = make(chan error) // Create a default channel if none provided
	}
	return &stage
}

func (s *ProcessorStage[IN, OUT]) OnInput(
	input IN,
) {
	// Process the input data and return the output
	if out, err := s.Process(input); err != nil {
		s.Error <- err // Send error to the error channel
	} else {
		if s.Output != nil {
			s.Output <- &out // Send output to the output channel
		}
	}
}

type PipelineStage[IN any, OUT any] struct {
	Name      string                     // Name of the stage, used for identification and logging
	Arguments <-chan IN                  // Input channel for the stage, where data is received
	Output    chan<- *OUT                // Output channel for the stage, where processed data is sent
	Error     chan<- error               // Error channel for the stage, where errors are sent
	Next      *[]PipelineStage[OUT, any] // Pointer to the next stage in the pipeline, allowing chaining of stages
}

func (s *PipelineStage[IN, OUT]) process(data IN) (OUT, error) {
	panic("unimplemented")
}

type PipelineStageInterface[IN any, OUT any] interface {
	StageName() string
	process(input IN) (OUT, error) // Process method to be implemented by the stage
}

func (s *PipelineStage[IN, OUT]) StageName() string {
	return s.Name
}
