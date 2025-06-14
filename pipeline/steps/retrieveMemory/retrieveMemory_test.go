package retrievememory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teagan42/snidemind/models"
)

func TestRetrieveMemory_Process_ReturnsInputUnchanged(t *testing.T) {
	step := &RetrieveMemory{}

	previous := []models.PipelineStep{}
	input := &models.PipelineMessage{
		Request: &models.ChatCompletionRequest{
			Messages: []models.ChatMessage{
				{
					Role:    "user",
					Content: "Hello, how are you?",
				},
			},
			Model: "gpt-3.5-turbo",
		},
	}

	output, err := step.Process(&previous, input)
	assert.NoError(t, err, "Process should not return an error")
	assert.Equal(t, input, output, "Process should return the input unchanged")
}

func TestRetrieveMemory_Process_NilInput(t *testing.T) {
	step := &RetrieveMemory{}

	previous := []models.PipelineStep{}
	var input *models.PipelineMessage = nil

	output, err := step.Process(&previous, input)
	assert.NoError(t, err, "Process should not return an error when input is nil")
	assert.Nil(t, output, "Process should return nil when input is nil")
}
