package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestChatCompletionRequest_JSONMarshalling(t *testing.T) {
	freqPenalty := 0.5
	presPenalty := 0.2
	maxTokens := int64(128)
	n := 2
	parallel := true
	stream := false
	temp := 0.9
	topP := 0.8
	topLogProbs := 3

	tools := &[]Tool{
		{
			Type: "function",
			Function: ToolFunction{
				Name:        "testFunc",
				Description: "desc",
				Parameters:  ToolFunctionParameters{"param1": "value1"},
			},
		},
	}

	webSearchOptions := &WebSearchOptions{
		SearchContextSize: "large",
		UserLocation: &WebSearchUserLocation{
			Type: "approximate",
			Approximate: ApproximateLocation{
				City:     "Seattle",
				Country:  "USA",
				Region:   "WA",
				TimeZone: "PST",
			},
		},
	}

	req := ChatCompletionRequest{
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Model:               "gpt-4",
		FrequencyPenalty:    &freqPenalty,
		FunctionCall:        &FunctionCall{Name: "myFunc"},
		LogProbs:            true,
		MaxCompletionTokens: &maxTokens,
		N:                   &n,
		ParallelToolCalls:   &parallel,
		PresencePenalty:     &presPenalty,
		ReasoningEffort:     "high",
		ResponseFormat:      "json",
		Stream:              &stream,
		Temperature:         &temp,
		Tools:               tools,
		TopLogProbs:         topLogProbs,
		TopP:                &topP,
		User:                "user-123",
		WebSearchOptions:    webSearchOptions,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatCompletionRequest: %v", err)
	}

	var unmarshalled ChatCompletionRequest
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal ChatCompletionRequest: %v", err)
	}

	// Check required fields
	if len(unmarshalled.Messages) != 1 || unmarshalled.Messages[0].Content != "Hello" {
		t.Errorf("Messages not marshalled/unmarshalled correctly")
	}
	if unmarshalled.Model != "gpt-4" {
		t.Errorf("Model not marshalled/unmarshalled correctly")
	}
	if unmarshalled.FrequencyPenalty == nil || *unmarshalled.FrequencyPenalty != freqPenalty {
		t.Errorf("FrequencyPenalty not marshalled/unmarshalled correctly")
	}
	if unmarshalled.FunctionCall == nil || unmarshalled.FunctionCall.Name != "myFunc" {
		t.Errorf("FunctionCall not marshalled/unmarshalled correctly")
	}
	if !unmarshalled.LogProbs {
		t.Errorf("LogProbs not marshalled/unmarshalled correctly")
	}
	if unmarshalled.MaxCompletionTokens == nil || *unmarshalled.MaxCompletionTokens != maxTokens {
		t.Errorf("MaxCompletionTokens not marshalled/unmarshalled correctly")
	}
	if unmarshalled.N == nil || *unmarshalled.N != n {
		t.Errorf("N not marshalled/unmarshalled correctly")
	}
	if unmarshalled.ParallelToolCalls == nil || *unmarshalled.ParallelToolCalls != parallel {
		t.Errorf("ParallelToolCalls not marshalled/unmarshalled correctly")
	}
	if unmarshalled.PresencePenalty == nil || *unmarshalled.PresencePenalty != presPenalty {
		t.Errorf("PresencePenalty not marshalled/unmarshalled correctly")
	}
	if unmarshalled.ReasoningEffort != "high" {
		t.Errorf("ReasoningEffort not marshalled/unmarshalled correctly")
	}
	if unmarshalled.ResponseFormat != "json" {
		t.Errorf("ResponseFormat not marshalled/unmarshalled correctly")
	}
	if unmarshalled.Stream == nil || *unmarshalled.Stream != stream {
		t.Errorf("Stream not marshalled/unmarshalled correctly")
	}
	if unmarshalled.Temperature == nil || *unmarshalled.Temperature != temp {
		t.Errorf("Temperature not marshalled/unmarshalled correctly")
	}
	if unmarshalled.Tools == nil || !reflect.DeepEqual(unmarshalled.Tools, tools) {
		t.Errorf("Tools not marshalled/unmarshalled correctly")
	}
	if unmarshalled.TopLogProbs != topLogProbs {
		t.Errorf("TopLogProbs not marshalled/unmarshalled correctly")
	}
	if unmarshalled.TopP == nil || *unmarshalled.TopP != topP {
		t.Errorf("TopP not marshalled/unmarshalled correctly")
	}
	if unmarshalled.User != "user-123" {
		t.Errorf("User not marshalled/unmarshalled correctly")
	}
	if unmarshalled.WebSearchOptions == nil || !reflect.DeepEqual(unmarshalled.WebSearchOptions, webSearchOptions) {
		t.Errorf("WebSearchOptions not marshalled/unmarshalled correctly")
	}
}

func TestChatCompletionRequest_EmptyOptionalFields(t *testing.T) {
	req := ChatCompletionRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: "Test"},
		},
		Model: "gpt-3.5",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatCompletionRequest: %v", err)
	}

	var unmarshalled ChatCompletionRequest
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal ChatCompletionRequest: %v", err)
	}

	if len(unmarshalled.Messages) != 1 || unmarshalled.Messages[0].Content != "Test" {
		t.Errorf("Messages not marshalled/unmarshalled correctly")
	}
	if unmarshalled.Model != "gpt-3.5" {
		t.Errorf("Model not marshalled/unmarshalled correctly")
	}
	// All optional fields should be nil or zero value
	if unmarshalled.FrequencyPenalty != nil {
		t.Errorf("Expected FrequencyPenalty to be nil")
	}
	if unmarshalled.FunctionCall != nil {
		t.Errorf("Expected FunctionCall to be nil")
	}
	if unmarshalled.LogProbs {
		t.Errorf("Expected LogProbs to be false")
	}
	if unmarshalled.MaxCompletionTokens != nil {
		t.Errorf("Expected MaxCompletionTokens to be nil")
	}
	if unmarshalled.N != nil {
		t.Errorf("Expected N to be nil")
	}
	if unmarshalled.ParallelToolCalls != nil {
		t.Errorf("Expected ParallelToolCalls to be nil")
	}
	if unmarshalled.PresencePenalty != nil {
		t.Errorf("Expected PresencePenalty to be nil")
	}
	if unmarshalled.ReasoningEffort != "" {
		t.Errorf("Expected ReasoningEffort to be empty")
	}
	if unmarshalled.ResponseFormat != "" {
		t.Errorf("Expected ResponseFormat to be empty")
	}
	if unmarshalled.Stream != nil {
		t.Errorf("Expected Stream to be nil")
	}
	if unmarshalled.Temperature != nil {
		t.Errorf("Expected Temperature to be nil")
	}
	if unmarshalled.Tools != nil {
		t.Errorf("Expected Tools to be nil")
	}
	if unmarshalled.TopLogProbs != 0 {
		t.Errorf("Expected TopLogProbs to be zero")
	}
	if unmarshalled.TopP != nil {
		t.Errorf("Expected TopP to be nil")
	}
	if unmarshalled.User != "" {
		t.Errorf("Expected User to be empty")
	}
	if unmarshalled.WebSearchOptions != nil {
		t.Errorf("Expected WebSearchOptions to be nil")
	}
}
