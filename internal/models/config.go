package models

import (
	"fmt"

	"github.com/teagan42/snidemind/internal/types"
)

type Config struct {
	Server     ServerConfig      `json:"server" yaml:"server" validate:"required"`
	MCPServers []MCPServerConfig `json:"mcp_servers" yaml:"mcp_servers" validate:"required"`
	LLM        LLMConfig         `json:"llm" yaml:"llm" validate:"required"`
}

type LLMConfig struct {
	Model            *string           `json:"model,omitempty" yaml:"model,omitempty" validate:"omitempty,required"`
	APIKey           *string           `json:"api_key,omitempty" yaml:"api_key,omitempty" validate:"omitempty,required"`
	APIKeyHeader     *string           `json:"api_key_header,omitempty" yaml:"api_key,omitempty" validate:"omitempty,required=api_key"`
	BaseURL          string            `json:"base_url,omitempty" yaml:"base_url,omitempty" validate:"omitempty,url"`
	Timeout          *int              `json:"timeout,omitempty" yaml:"timeout,omitempty" validate:"omitempty,min=1"`
	Headers          map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" validate:"omitempty,dive,keys,required"`
	Temperature      *float64          `json:"temperature,omitempty" yaml:"temperature,omitempty" validate:"omitempty,min=0,max=1"`
	MaxTokens        *int64            `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty" validate:"omitempty,min=1"`
	TopP             *float64          `json:"top_p,omitempty" yaml:"top_p,omitempty" validate:"omitempty,min=0,max=1"`
	FrequencyPenalty *float64          `json:"frequency_penalty,omitempty" yaml:"frequency_penalty,omitempty" validate:"omitempty,min=0,max=1"`
	PresencePenalty  *float64          `json:"presence_penalty,omitempty" yaml:"presence_penalty,omitempty" validate:"omitempty,min=0,max=1"`
	N                *int              `json:"n,omitempty" yaml:"n,omitempty" validate:"omitempty,min=1"`
	Stream           *bool             `json:"stream,omitempty" yaml:"stream,omitempty" validate:"omitempty"`
}

type MCPBlacklist struct {
	Tools     *types.RegexList `json:"tools,omitempty" yaml:"tools,omitempty" validate:"omitempty,dive,required"`
	Prompts   *types.RegexList `json:"prompts,omitempty"  yaml:"prompts,omitempty" validate:"omitempty,dive,required"`
	Resources *types.RegexList `json:"resources,omitempty" yaml:"resources,omitempty" validate:"omitempty,dive,required"`
}

type MCPServerConfig struct {
	Name      string        `json:"name" yaml:"name" validate:"required"`
	URL       string        `json:"url" yaml:"url" validate:"required,url"`
	Type      string        `json:"type" yaml:"type" validate:"required,oneof=sse http"`
	Blacklist *MCPBlacklist `json:"blacklist,omitempty" yaml:"blacklist,omitempty" validate:"omitempty,dive"`
}

type ServerConfig struct {
	Port types.Port `json:"port" yaml:"port" validate:"required"`
	Bind types.Host `json:"bind" yaml:"bind" validate:"required"`
}

func (b *MCPBlacklist) Validate() error {
	for _, lst := range []*types.RegexList{b.Tools, b.Prompts, b.Resources} {
		if lst == nil {
			continue
		}
		for _, re := range *lst {
			if re == nil {
				return fmt.Errorf("nil regex found in blacklist")
			}
		}
	}
	return nil
}

func (b *MCPBlacklist) IsToolBlacklisted(toolName string) bool {
	if b == nil || b.Tools == nil {
		return false
	}
	for _, t := range *b.Tools {
		if t.Match([]byte(toolName)) {
			return true
		}

	}
	return false
}

func (b *MCPBlacklist) IsPromptBlacklisted(promptName string) bool {
	if b == nil || b.Tools == nil {
		return false
	}
	for _, t := range *b.Tools {
		if t.Match([]byte(promptName)) {
			return true
		}

	}
	return false
}

func (b *MCPBlacklist) IsResourceBlacklisted(resourceName string) bool {
	if b == nil || b.Resources == nil {
		return false
	}
	for _, t := range *b.Resources {
		if t.Match([]byte(resourceName)) {
			return true
		}

	}
	return false
}
