// File: config/types_test.go
package config

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestLLMConfig_Validation(t *testing.T) {
	model := "gpt-4"
	apiKey := "sk-123"
	apiKeyHeader := "Authorization"
	baseURL := "https://api.example.com"
	timeout := 30
	headers := map[string]string{"X-Test": "value"}
	temperature := 0.7
	maxTokens := int64(256)
	topP := 0.9
	frequencyPenalty := 0.1
	presencePenalty := 0.2
	n := 2
	stream := true
	parallelToolCalls := true

	tests := []struct {
		name  string
		cfg   LLMConfig
		valid bool
	}{
		{
			name: "All fields set and valid",
			cfg: LLMConfig{
				Model:             &model,
				APIKey:            &apiKey,
				APIKeyHeader:      &apiKeyHeader,
				BaseURL:           baseURL,
				Timeout:           &timeout,
				Headers:           headers,
				Temperature:       &temperature,
				MaxTokens:         &maxTokens,
				TopP:              &topP,
				FrequencyPenalty:  &frequencyPenalty,
				PresencePenalty:   &presencePenalty,
				N:                 &n,
				Stream:            &stream,
				ParallelToolCalls: &parallelToolCalls,
			},
			valid: true,
		},
		{
			name: "Invalid BaseURL",
			cfg: LLMConfig{
				Model:   &model,
				APIKey:  &apiKey,
				BaseURL: "not-a-url",
			},
			valid: false,
		},
		{
			name: "Invalid Temperature (too high)",
			cfg: LLMConfig{
				Model:       &model,
				APIKey:      &apiKey,
				Temperature: floatPtr(2.0),
			},
			valid: false,
		},
		{
			name: "Invalid Temperature (too low)",
			cfg: LLMConfig{
				Model:       &model,
				APIKey:      &apiKey,
				Temperature: floatPtr(-1.0),
			},
			valid: false,
		},
		{
			name: "Invalid MaxTokens (zero)",
			cfg: LLMConfig{
				Model:     &model,
				APIKey:    &apiKey,
				MaxTokens: int64Ptr(0),
			},
			valid: false,
		},
		{
			name: "Invalid TopP (negative)",
			cfg: LLMConfig{
				Model:  &model,
				APIKey: &apiKey,
				TopP:   floatPtr(-0.1),
			},
			valid: false,
		},
		{
			name: "Valid minimal config",
			cfg: LLMConfig{
				Model:  &model,
				APIKey: &apiKey,
			},
			valid: true,
		},
	}

	// Use github.com/go-playground/validator/v10 for validation
	validate := getValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.cfg)
			if tt.valid && err != nil {
				t.Errorf("expected valid config, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected invalid config, but got no error")
			}
		})
	}
}

func floatPtr(f float64) *float64 { return &f }
func int64Ptr(i int64) *int64     { return &i }

// getValidator returns a validator instance with required tag support.
func getValidator() *validator.Validate {
	// Import here for test only
	// go get github.com/go-playground/validator/v10
	v := validator.New()
	return v
}
