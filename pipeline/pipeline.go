package pipeline

import (
	"net/http"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"github.com/teagan42/snidemind/server/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Pipeline struct {
	Steps  []models.PipelineStep // All Steps in the pipeline
	Logger *zap.Logger
}

type Params struct {
	fx.In
	Config        *config.Config
	Logger        *zap.Logger
	StepFactories map[string]models.PipelineStepFactory `name:"pipelineStepFactoryMap"`
}

func NewPipeline(p Params) *Pipeline {
	steps := []models.PipelineStep{}
	for _, step := range p.Config.Pipeline.Steps {
		if factory, ok := p.StepFactories[step.Type]; ok {
			if s, err := factory.Build(step, p.StepFactories); err != nil {
				p.Logger.Error("Error building pipeline step", zap.Error(err))
				continue
			} else {
				steps = append(steps, s)
			}
		} else {
			p.Logger.Error("Unknown pipeline step type", zap.String("type", step.Type))
		}
	}
	p.Logger.Info("Pipeline initialized", zap.Int("steps", len(steps)))
	return &Pipeline{
		Steps:  steps,
		Logger: p.Logger.Named("Pipeline"),
	}
}

func (p *Pipeline) AddStep(stage models.PipelineStep, index *int) {
	if index == nil || *index < 0 || *index > len(p.Steps) {
		p.Steps = append(p.Steps, stage)
	} else {
		p.Steps = append(p.Steps[:*index], append([]models.PipelineStep{stage}, p.Steps[*index:]...)...)
	}
}

func (p *Pipeline) Process(r *http.Request, w http.ResponseWriter) (models.PipelineMessage, error) {
	p.Logger.Info("Processing pipeline for request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
	body, err := middleware.GetValidatedBody[models.ChatCompletionRequest](r)
	if err != nil {
		p.Logger.Error("Error getting validated body", zap.Error(err))
		return *new(models.PipelineMessage), err // Return zero value of OUT and the error
	}
	p.Logger.Info("Validated body", zap.Any("body", body))
	input := &models.PipelineMessage{
		Request:        &body,
		Tags:           &[]string{},
		Tools:          &[]string{},
		Prompts:        &[]string{},
		Memories:       &[]string{},
		Knowledge:      &[]string{},
		ResponseWriter: w,
	}
	var previous *[]models.PipelineStep
	for _, stage := range p.Steps {
		p.Logger.Info("Processing stage", zap.String("stage", stage.Name()))
		if input, err = stage.Process(previous, input); err != nil {
			p.Logger.Error("Error processing stage", zap.String("stage", stage.Name()), zap.Error(err))
			return *new(models.PipelineMessage), err // Return zero value of OUT and the error
		}
		p.Logger.Info("Stage processed successfully", zap.String("stage", stage.Name()))
		previous = &[]models.PipelineStep{stage} // Update previous to the current stage
	}

	return *input, nil // Return the processed output
}
