package fork

import (
	"fmt"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ForkedPipelineStages struct {
	Steps []models.PipelineStep
}

type ForkPipelineStage struct {
	Forks  []ForkedPipelineStages
	Logger *zap.Logger
}

type Params struct {
	fx.In
	Config *config.Config
	Logger *zap.Logger
}

type Result struct {
	fx.Out
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type ForkPipelineStageFactory struct {
	Logger *zap.Logger
}

func (f ForkPipelineStageFactory) Name() string {
	return "fork"
}

func (f ForkPipelineStageFactory) Build(config config.PipelineStepConfig, stepFactories map[string]models.PipelineStepFactory) (models.PipelineStep, error) {
	f.Logger.Info("Building Fork Stage", zap.Any("config", config))
	if config.Fork == nil {
		return nil, fmt.Errorf("fork config is nil")
	}
	if len(*config.Fork) == 0 {
		return nil, fmt.Errorf("fork config is empty")
	}
	var forkedStages = []ForkedPipelineStages{}
	for _, forkConfig := range *config.Fork {
		f.Logger.Info("Processing Fork Config", zap.Any("forkConfig", forkConfig))
		var steps = []models.PipelineStep{}
		for j, step := range forkConfig.Steps {
			f.Logger.Info("Processing Step", zap.Int("index", j), zap.Any("step", step))
			if step.Type == "" {
				f.Logger.Error("Step type is empty in fork config", zap.Int("index", j))
				return nil, fmt.Errorf("step type is empty in fork config at index %d", j)
			}
			factory, ok := stepFactories[step.Type]
			if !ok {
				f.Logger.Error("Unknown pipeline step type", zap.String("type", step.Type))
				return nil, fmt.Errorf("unknown pipeline step type: %s", step.Type)
			}
			stage, err := factory.Build(step, stepFactories)
			if err != nil {
				f.Logger.Error("Error building pipeline step", zap.String("type", step.Type), zap.Error(err))
				return nil, fmt.Errorf("failed to build pipeline step: %w", err)
			}
			steps = append(steps, stage)
		}
		f.Logger.Info("Forked Steps", zap.Int("count", len(steps)), zap.Any("steps", steps))
		forkedStages = append(forkedStages, ForkedPipelineStages{Steps: steps})
	}
	f.Logger.Info("Fork Stage Built", zap.Int("forks", len(forkedStages)), zap.Any("forkedStages", forkedStages))
	return &ForkPipelineStage{
		Forks:  forkedStages,
		Logger: f.Logger.Named("ForkPipelineStage"),
	}, nil
}

func NewFork(p Params) (Result, error) {
	return Result{
		Factory: ForkPipelineStageFactory{
			Logger: p.Logger.Named("ForkPipelineStage"),
		},
	}, nil
}

func (f ForkPipelineStage) Name() string {
	return "fork"
}

func (f ForkPipelineStage) processFork(previous *[]models.PipelineStep, steps []models.PipelineStep, input *models.PipelineMessage, resultsChannel chan *models.PipelineMessage) {
	f.Logger.Info("Processing Forked Steps", zap.Any("steps", steps))
	for _, step := range steps {
		var err error
		input, err = step.Process(previous, input)
		if err != nil {
			fmt.Printf("Error processing step: %v\n", err)
			resultsChannel <- nil
			return
		}
		previous = &[]models.PipelineStep{step}
	}
	resultsChannel <- input
}

func (f ForkPipelineStage) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	f.Logger.Info("Processing Fork Stage")
	var resultsChannel = make(chan *models.PipelineMessage, len(f.Forks))
	defer close(resultsChannel)
	for _, fork := range f.Forks {
		go f.processFork(previous, fork.Steps, input, resultsChannel)
	}

	resultsDone := 0
	for result := range resultsChannel {
		resultsDone++
		if result != nil {
			input.Combine(result)
		}
		if resultsDone == len(f.Forks) {
			break
		}
	}
	return input, nil
}
