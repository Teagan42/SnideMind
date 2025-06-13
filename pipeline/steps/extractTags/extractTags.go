package extracttags

import (
	"fmt"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ExtractTags struct {
	Embedder *Embedder
}

type Params struct {
	fx.In
	Log *zap.Logger
}

type Result struct {
	fx.Out
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type ExtractTagsFactory struct {
	Logger    *zap.Logger
	Embedders map[string]*Embedder
}

func (f ExtractTagsFactory) Name() string {
	return "extractTags"
}
func (f ExtractTagsFactory) Build(config config.PipelineStepConfig, stepFactories map[string]models.PipelineStepFactory) (models.PipelineStep, error) {
	if config.Embedder == nil {
		return nil, fmt.Errorf("extractTags config is nil")
	}
	stepConfig := *config.Embedder
	key := stepConfig.URL + stepConfig.Model
	if _, ok := f.Embedders[key]; !ok {
		f.Embedders[key] = NewEmbedder(
			f.Logger,
			stepConfig.URL,
			stepConfig.Model,
			stepConfig.APIKey,
			stepConfig.APIKeyHeader,
		)
	}
	return &ExtractTags{
		Embedder: f.Embedders[key],
	}, nil
}

func NewExtractTags(p Params) (Result, error) {
	return Result{
		Factory: ExtractTagsFactory{
			Logger:    p.Log,
			Embedders: map[string]*Embedder{},
		},
	}, nil
}

func (s ExtractTags) Name() string {
	return "extractTags"
}

func (s ExtractTags) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	if input.Request == nil {
		return nil, fmt.Errorf("input request is nil")
	}

	if input.Request.Messages == nil {
		return nil, fmt.Errorf("input messages are nil")
	}
	msg := input.Request.Messages[len(input.Request.Messages)-1].Content

	if tags, error := s.Embedder.ExtractTagsWithWeights(msg); error != nil {
		return nil, fmt.Errorf("failed to extract tags: %w", error)
	} else {
		s.Embedder.Logger.Info("Extracted tags", zap.String("tags", fmt.Sprintf("%v", tags)))
		if input.Tags == nil {
			input.Tags = &map[string]string{}
		}
		for _, tag := range tags {
			(*input.Tags)[tag.Tag.ID] = tag.Tag.ID
		}
	}
	return input, nil
}
