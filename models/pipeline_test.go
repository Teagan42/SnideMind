// File: models/pipeline_test.go
package models

import (
	"reflect"
	"testing"
)

func TestPipelineMessage_Combine_AllFields(t *testing.T) {
	// Prepare initial PipelineMessage
	pm1 := &PipelineMessage{
		Tags: &map[string]string{"a": "1"},
		Tools: &[]MCPTool{{
			ToolMetadata{
				Name:        "tool1",
				Description: "description1",
				Tags:        &[]string{"tag1", "tag2"},
			},
		}},
		Prompts:   &[]string{"prompt1"},
		Memories:  &[]string{"memory1"},
		Knowledge: &[]string{"knowledge1"},
	}

	// Prepare message to combine
	pm2 := &PipelineMessage{
		Tags: &map[string]string{"b": "2"},
		Tools: &[]MCPTool{{
			ToolMetadata{
				Name:        "tool2",
				Description: "description2",
				Tags:        &[]string{"tag3", "tag4"},
			},
		}},
		Prompts:   &[]string{"prompt2"},
		Memories:  &[]string{"memory2"},
		Knowledge: &[]string{"knowledge2"},
	}

	pm1.Combine(pm2)

	// Check Tags
	wantTags := map[string]string{"a": "1", "b": "2"}
	if pm1.Tags == nil || !reflect.DeepEqual(*pm1.Tags, wantTags) {
		t.Errorf("Tags not combined correctly: got %v, want %v", pm1.Tags, wantTags)
	}

	// Check Tools
	wantTools := []MCPTool{
		{
			ToolMetadata{
				Name:        "tool1",
				Description: "description1",
				Tags:        &[]string{"tag1", "tag2"},
			},
		},
		{
			ToolMetadata{
				Name:        "tool2",
				Description: "description2",
				Tags:        &[]string{"tag3", "tag4"},
			},
		},
	}
	if pm1.Tools == nil || !reflect.DeepEqual(*pm1.Tools, wantTools) {
		t.Errorf("Tools not combined correctly: got %v, want %v", pm1.Tools, wantTools)
	}

	// Check Prompts
	wantPrompts := []string{"prompt1", "prompt2"}
	if pm1.Prompts == nil || !reflect.DeepEqual(*pm1.Prompts, wantPrompts) {
		t.Errorf("Prompts not combined correctly: got %v, want %v", pm1.Prompts, wantPrompts)
	}

	// Check Memories
	wantMemories := []string{"memory1", "memory2"}
	if pm1.Memories == nil || !reflect.DeepEqual(*pm1.Memories, wantMemories) {
		t.Errorf("Memories not combined correctly: got %v, want %v", pm1.Memories, wantMemories)
	}

	// Check Knowledge
	wantKnowledge := []string{"knowledge1", "knowledge2"}
	if pm1.Knowledge == nil || !reflect.DeepEqual(*pm1.Knowledge, wantKnowledge) {
		t.Errorf("Knowledge not combined correctly: got %v, want %v", pm1.Knowledge, wantKnowledge)
	}
}

func TestPipelineMessage_Combine_NilFields(t *testing.T) {
	pm1 := &PipelineMessage{}
	pm2 := &PipelineMessage{
		Tags: &map[string]string{"x": "y"},
		Tools: &[]MCPTool{{
			ToolMetadata{
				Name:        "t",
				Description: "d",
				Tags:        &[]string{"tag1", "tag2"},
			},
		}},
		Prompts:   &[]string{"p"},
		Memories:  &[]string{"m"},
		Knowledge: &[]string{"k"},
	}

	pm1.Combine(pm2)

	if pm1.Tags == nil || (*pm1.Tags)["x"] != "y" {
		t.Errorf("Tags not set correctly when initial is nil")
	}
	if pm1.Tools == nil || len(*pm1.Tools) != 1 || (*pm1.Tools)[0].ToolMetadata.Name != "t" {
		t.Errorf("Tools not set correctly when initial is nil")
	}
	if pm1.Prompts == nil || len(*pm1.Prompts) != 1 || (*pm1.Prompts)[0] != "p" {
		t.Errorf("Prompts not set correctly when initial is nil")
	}
	if pm1.Memories == nil || len(*pm1.Memories) != 1 || (*pm1.Memories)[0] != "m" {
		t.Errorf("Memories not set correctly when initial is nil")
	}
	if pm1.Knowledge == nil || len(*pm1.Knowledge) != 1 || (*pm1.Knowledge)[0] != "k" {
		t.Errorf("Knowledge not set correctly when initial is nil")
	}
}

func TestPipelineMessage_Combine_NilInputFields(t *testing.T) {
	pm1 := &PipelineMessage{
		Tags: &map[string]string{"a": "1"},
		Tools: &[]MCPTool{{
			ToolMetadata{
				Name:        "tool1",
				Description: "description1",
				Tags:        &[]string{"tag1", "tag2"},
			},
		}},
		Prompts:   &[]string{"prompt1"},
		Memories:  &[]string{"memory1"},
		Knowledge: &[]string{"knowledge1"},
	}
	pm2 := &PipelineMessage{} // All fields nil

	orig := *pm1

	pm1.Combine(pm2)

	// Should remain unchanged
	if !reflect.DeepEqual(*pm1.Tags, *orig.Tags) {
		t.Errorf("Tags changed unexpectedly: got %v, want %v", pm1.Tags, orig.Tags)
	}
	if !reflect.DeepEqual(*pm1.Tools, *orig.Tools) {
		t.Errorf("Tools changed unexpectedly: got %v, want %v", pm1.Tools, orig.Tools)
	}
	if !reflect.DeepEqual(*pm1.Prompts, *orig.Prompts) {
		t.Errorf("Prompts changed unexpectedly: got %v, want %v", pm1.Prompts, orig.Prompts)
	}
	if !reflect.DeepEqual(*pm1.Memories, *orig.Memories) {
		t.Errorf("Memories changed unexpectedly: got %v, want %v", pm1.Memories, orig.Memories)
	}
	if !reflect.DeepEqual(*pm1.Knowledge, *orig.Knowledge) {
		t.Errorf("Knowledge changed unexpectedly: got %v, want %v", pm1.Knowledge, orig.Knowledge)
	}
}
