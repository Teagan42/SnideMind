package models

type ToolArgument struct {
	Name        string `json:"name" validate:"required"`                                           // Name of the argument
	Description string `json:"description,omitempty"`                                              // Optional description of the argument
	Type        string `json:"type" validate:"required,oneof=string integer boolean array object"` // Type of the argument
}

type Tool struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=tool"` // "function" for OpenAI, "tool" for Claude
	Function string `json:"function" validate:"required"`        // Function name or description
	Tag      string `json:"tag,omitempty"`                       // Optional tag for categorization
}

type Tools struct {
	Tools []Tool `json:"tools" validate:"required,dive"` // List of tools
}

type ErrorInvalidTool struct {
	Err string `json:"error" validate:"required"` // Error message
}

func (e ErrorInvalidTool) Error() string {
	return e.Err
}

func NewTool(tool Tool) (Tool, error) {
	if tool.ID == "" || tool.Name == "" || tool.Type == "" || tool.Function == "" {
		return Tool{}, ErrorInvalidTool{
			"Tool must have ID, Name, Type, and Function fields set",
		}
	}
	if tool.Type != "tool" && tool.Type != "function" {
		return Tool{}, ErrorInvalidTool{
			"Tool Type must be either 'tool' or 'function'",
		}
	}
	return tool, nil
}
