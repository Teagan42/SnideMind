package reducetools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/zap"
)

func TestNewReduceTools(t *testing.T) {
	logger := zap.NewNop()
	params := Params{
		Logger: logger,
	}

	result, err := NewReduceTools(params)
	assert.NoError(t, err, "NewReduceTools should not return an error")
	assert.NotNil(t, result.Factory, "Factory should not be nil")

	factory, ok := result.Factory.(ReduceToolsFactory)
	assert.True(t, ok, "Factory should be of type ReduceToolsFactory")
	assert.NotNil(t, factory.Logger, "Logger should not be nil")
	assert.Equal(t, "ReduceToolsFactory", factory.Logger.Name(), "Logger should be named ReduceToolsFactory")

	tools := factory.ToolSet
	assert.Len(t, tools, 3, "ToolSet should contain 3 tools")

	toolNames := []string{
		"Plex New Media",
		"Plex Search",
		"Home Assistant Entity Action",
	}
	for i, tool := range tools {
		assert.Equal(t, toolNames[i], tool.ToolMetadata.Name)
		assert.NotEmpty(t, tool.ToolMetadata.Description)
		assert.NotNil(t, tool.ToolMetadata.Tags)
		assert.Greater(t, len(*tool.ToolMetadata.Tags), 0)
	}
}
func TestReduceTools_Process_NoTags(t *testing.T) {
	logger := zap.NewNop()
	toolSet := []models.MCPTool{
		{
			ToolMetadata: models.ToolMetadata{
				Name:        "Tool1",
				Description: "Desc1",
				Tags:        &[]string{"tag1", "tag2"},
			},
		},
	}
	rt := ReduceTools{
		Logger:  logger,
		ToolSet: toolSet,
	}
	input := &models.PipelineMessage{
		Tags: nil,
	}
	prev := &[]models.PipelineStep{}
	result, err := rt.Process(prev, input)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
	assert.Nil(t, result.Tools)
}

func TestReduceTools_Process_EmptyTags(t *testing.T) {
	logger := zap.NewNop()
	toolSet := []models.MCPTool{
		{
			ToolMetadata: models.ToolMetadata{
				Name:        "Tool1",
				Description: "Desc1",
				Tags:        &[]string{"tag1", "tag2"},
			},
		},
	}
	rt := ReduceTools{
		Logger:  logger,
		ToolSet: toolSet,
	}
	input := &models.PipelineMessage{
		Tags: &map[string]string{},
	}
	prev := &[]models.PipelineStep{}
	result, err := rt.Process(prev, input)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
	assert.Nil(t, result.Tools)
}

func TestReduceTools_Process_MatchingTags(t *testing.T) {
	logger := zap.NewNop()
	tool1 := models.MCPTool{
		ToolMetadata: models.ToolMetadata{
			Name:        "Tool1",
			Description: "Desc1",
			Tags:        &[]string{"tag1", "tag2"},
		},
	}
	tool2 := models.MCPTool{
		ToolMetadata: models.ToolMetadata{
			Name:        "Tool2",
			Description: "Desc2",
			Tags:        &[]string{"tag3"},
		},
	}
	rt := ReduceTools{
		Logger:  logger,
		ToolSet: []models.MCPTool{tool1, tool2},
	}
	tags := map[string]string{"tag1": "v", "tag3": "v"}
	input := &models.PipelineMessage{
		Tags: &tags,
	}
	prev := &[]models.PipelineStep{}
	result, err := rt.Process(prev, input)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
	assert.NotNil(t, result.Tools)
	assert.Len(t, *result.Tools, 2)
	names := []string{(*result.Tools)[0].ToolMetadata.Name, (*result.Tools)[1].ToolMetadata.Name}
	assert.Contains(t, names, "Tool1")
	assert.Contains(t, names, "Tool2")
}

func TestReduceTools_Process_PartialMatch(t *testing.T) {
	logger := zap.NewNop()
	tool1 := models.MCPTool{
		ToolMetadata: models.ToolMetadata{
			Name:        "Tool1",
			Description: "Desc1",
			Tags:        &[]string{"tag1", "tag2"},
		},
	}
	tool2 := models.MCPTool{
		ToolMetadata: models.ToolMetadata{
			Name:        "Tool2",
			Description: "Desc2",
			Tags:        &[]string{"tag3"},
		},
	}
	rt := ReduceTools{
		Logger:  logger,
		ToolSet: []models.MCPTool{tool1, tool2},
	}
	tags := map[string]string{"tag2": "v"}
	input := &models.PipelineMessage{
		Tags: &tags,
	}
	prev := &[]models.PipelineStep{}
	result, err := rt.Process(prev, input)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
	assert.NotNil(t, result.Tools)
	assert.Len(t, *result.Tools, 1)
	assert.Equal(t, "Tool1", (*result.Tools)[0].ToolMetadata.Name)
}

func TestReduceTools_Process_NoMatchingTags(t *testing.T) {
	logger := zap.NewNop()
	tool1 := models.MCPTool{
		ToolMetadata: models.ToolMetadata{
			Name:        "Tool1",
			Description: "Desc1",
			Tags:        &[]string{"tag1", "tag2"},
		},
	}
	rt := ReduceTools{
		Logger:  logger,
		ToolSet: []models.MCPTool{tool1},
	}
	tags := map[string]string{"tagX": "v"}
	input := &models.PipelineMessage{
		Tags: &tags,
	}
	prev := &[]models.PipelineStep{}
	result, err := rt.Process(prev, input)
	assert.NoError(t, err)
	assert.Equal(t, input, result)
	assert.NotNil(t, result.Tools)
	assert.Len(t, *result.Tools, 0)
}
