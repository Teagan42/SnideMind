package extracttags

import (
	"github.com/teagan42/snidemind/pipeline"
)

type ExtractTags struct {
}

func NewExtractTags() *ExtractTags {
	return &ExtractTags{}
}

func (s *ExtractTags) Name() string {
	return "ExtractTags"
}

func (s *ExtractTags) Process(input *pipeline.PipelineMessage) (*pipeline.PipelineMessage, error) {
	// Placeholder for tag extraction logic
	// In a real implementation, this would parse the input string and extract tags
	// For now, we'll just return a dummy slice of tags
	input.Tags = &[]string{"tag1", "tag2", "tag3"}
	return input, nil
}
