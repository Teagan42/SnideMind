package extracttags

import (
	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
)

type ExtractTags struct {
}

type Params struct {
	fx.In
	Config config.PipelineStepConfig
}

type Result struct {
	fx.Out
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type ExtractTagsFactory struct{}

func (f ExtractTagsFactory) Name() string {
	return "extractTags"
}
func (f ExtractTagsFactory) Build(config config.PipelineStepConfig) (models.PipelineStep, error) {
	return &ExtractTags{}, nil
}

func NewExtractTags() (Result, error) {
	return Result{
		Factory: ExtractTagsFactory{},
	}, nil
}

func (s ExtractTags) Name() string {
	return "ExtractTags"
}

func (s ExtractTags) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	// Placeholder for tag extraction logic
	// In a real implementation, this would parse the input string and extract tags
	// For now, we'll just return a dummy slice of tags
	input.Tags = &[]string{"tag1", "tag2", "tag3"}
	return input, nil
}
