package pipeline

type PipelineBuilder[IN any, OUT any] struct {
	Stages []PipelineStage[IN, OUT] // List of stages in the pipeline
}
