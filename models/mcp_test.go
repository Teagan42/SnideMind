// File: models/mcp_test.go
package models

import (
	"encoding/json"
	"testing"
)

func TestToolMetadata_JSONMarshalling(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	meta := ToolMetadata{
		Name:        "TestTool",
		Description: "A test tool",
		Tags:        &tags,
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("Failed to marshal ToolMetadata: %v", err)
	}

	var unmarshalled ToolMetadata
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal ToolMetadata: %v", err)
	}

	if unmarshalled.Name != meta.Name {
		t.Errorf("Expected Name %q, got %q", meta.Name, unmarshalled.Name)
	}
	if unmarshalled.Description != meta.Description {
		t.Errorf("Expected Description %q, got %q", meta.Description, unmarshalled.Description)
	}
	if unmarshalled.Tags == nil || len(*unmarshalled.Tags) != 2 {
		t.Errorf("Expected Tags length 2, got %v", unmarshalled.Tags)
	}
}

func TestToolMetadata_EmptyFields(t *testing.T) {
	meta := ToolMetadata{
		Name: "TestTool",
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("Failed to marshal ToolMetadata with empty fields: %v", err)
	}

	// Description and Tags should be omitted if empty/nil
	if string(data) != `{"name":"TestTool"}` {
		t.Errorf("Expected JSON to only include name, got: %s", data)
	}
}

func TestToolMetadata_RequiredName(t *testing.T) {
	meta := ToolMetadata{}
	if meta.Name != "" {
		t.Errorf("Expected Name to be empty, got %q", meta.Name)
	}
}
