package models

type Config struct {
	Server     ServerConfig      `json:"server" yaml:"server" validate:"required,dive"`
	MCPServers []MCPServerConfig `json:"mcp_servers" yaml:"mcp_servers" validate:"required,dive"`
}
