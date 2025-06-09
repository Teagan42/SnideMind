package models

import (
	"fmt"

	"github.com/teagan42/snidemind/internal/types"
)

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
