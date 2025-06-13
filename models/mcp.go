package models

type ToolMetadata struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

type MCPTool struct {
	ToolMetadata ToolMetadata `json:"metadata" validate:"required,dive"`
}
