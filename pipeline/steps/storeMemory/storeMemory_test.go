package storememory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teagan42/snidemind/models"
)

func TestStoreMemory_Process_ReturnsInputUnchanged(t *testing.T) {
	storeMemory := StoreMemory{}
	prevSteps := []models.PipelineStep{}
	input := &models.PipelineMessage{
		Request: &models.ChatCompletionRequest{
			Messages: []models.ChatMessage{
				{
					Role:    "user",
					Content: "Hello, world!",
				},
			},
			Model: "gpt-3.5-turbo",
		},
	}

	result, err := storeMemory.Process(&prevSteps, input)

	assert.NoError(t, err, "Process should not return an error")
	assert.Equal(t, input, result, "Process should return the input unchanged")
}

func TestStoreMemory_Process_NilInput(t *testing.T) {
	storeMemory := StoreMemory{}
	prevSteps := []models.PipelineStep{}

	result, err := storeMemory.Process(&prevSteps, nil)

	assert.NoError(t, err, "Process should not return an error when input is nil")
	assert.Nil(t, result, "Process should return nil when input is nil")
}
