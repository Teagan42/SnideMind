package config

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
)

type Host string

func isHostname(s string) bool {
	// A lazy hostname check. You can replace this with real DNS validation if you hate yourself.
	return len(s) > 0 && len(s) <= 255
}

func NewHost(val string) (Host, error) {
	if net.ParseIP(val) == nil && !isHostname(val) {
		return "", fmt.Errorf("invalid bind address: '%s'", val)
	}
	return Host(val), nil
}

type Port int

func NewPort(val int) (Port, error) {
	if val <= 0 || val > 65535 {
		return 0, fmt.Errorf("invalid port: %d", val)
	}
	return Port(val), nil
}

type RegexList []*regexp.Regexp

func (r *RegexList) UnmarshalJSON(data []byte) error {
	var raw []string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, pattern := range raw {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex %q: %w", pattern, err)
		}
		*r = append(*r, re)
	}

	return nil
}

func (r *RegexList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw []string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	for _, pattern := range raw {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex %q: %w", pattern, err)
		}
		*r = append(*r, re)
	}

	return nil
}

func (r RegexList) MarshalJSON() ([]byte, error) {
	var raw []string
	for _, re := range r {
		raw = append(raw, re.String())
	}
	return json.Marshal(raw)
}

func (r RegexList) MarshalYAML() (interface{}, error) {
	var raw []string
	for _, re := range r {
		raw = append(raw, re.String())
	}
	return raw, nil
}

type Config struct {
	Server     ServerConfig      `json:"server" yaml:"server" validate:"required"`
	MCPServers []MCPServerConfig `json:"mcp_servers" yaml:"mcp_servers" validate:"required"`
	Pipeline   PipelineConfig    `json:"pipeline,omitempty" yaml:"pipeline,omitempty" validate:"omitempty"`
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

type PipelineStepConfig struct {
	Type string            `json:"type" yaml:"type" validate:"required,oneof=extractTags fork llm reduceTools retrieveMemory storeMemory"`
	LLM  *LLMConfig        `json:"llm,omitempty" yaml:"llm,omitempty" validate:"omitempty"`
	Fork *[]PipelineConfig `json:"fork,omitempty" yaml:"fork,omitempty" validate:"omitempty"`
}

type PipelineConfig struct {
	Steps []PipelineStepConfig `json:"steps,omitempty" yaml:"steps,omitempty" validate:"omitempty,dive"`
}

type MCPBlacklist struct {
	Tools     *RegexList `json:"tools,omitempty" yaml:"tools,omitempty" validate:"omitempty,dive,required"`
	Prompts   *RegexList `json:"prompts,omitempty"  yaml:"prompts,omitempty" validate:"omitempty,dive,required"`
	Resources *RegexList `json:"resources,omitempty" yaml:"resources,omitempty" validate:"omitempty,dive,required"`
}

type MCPServerConfig struct {
	Name      string        `json:"name" yaml:"name" validate:"required"`
	URL       string        `json:"url" yaml:"url" validate:"required,url"`
	Type      string        `json:"type" yaml:"type" validate:"required,oneof=sse http"`
	Blacklist *MCPBlacklist `json:"blacklist,omitempty" yaml:"blacklist,omitempty" validate:"omitempty,dive"`
}

type ServerConfig struct {
	Port Port  `json:"port" yaml:"port" validate:"required"`
	Bind *Host `json:"bind" yaml:"bind" validate:"omitempty"`
}

func (b *MCPBlacklist) Validate() error {
	for _, lst := range []*RegexList{b.Tools, b.Prompts, b.Resources} {
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
