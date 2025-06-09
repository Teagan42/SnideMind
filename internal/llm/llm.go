package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/teagan42/snidemind/internal/models"
	"github.com/teagan42/snidemind/internal/types"
)

type LLM struct {
	Config models.LLMConfig
}

type LLMError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e LLMError) Error() string {
	return e.Message
}

func NewLLM(cfg models.LLMConfig) *LLM {
	return &LLM{
		Config: cfg,
	}
}

func (llm *LLM) BuildCompletionRequest(request types.ChatCompletionRequest) types.ChatCompletionRequest {
	body := types.ChatCompletionRequest{
		Messages: request.Messages,
		Model:    request.Model,
	}
	if request.FrequencyPenalty != nil {
		body.FrequencyPenalty = request.FrequencyPenalty
	} else if llm.Config.FrequencyPenalty != nil {
		body.FrequencyPenalty = llm.Config.FrequencyPenalty
	}

	if request.PresencePenalty != nil {
		body.PresencePenalty = request.PresencePenalty
	} else if llm.Config.PresencePenalty != nil {
		body.PresencePenalty = llm.Config.PresencePenalty
	}
	if request.Temperature != nil {
		body.Temperature = request.Temperature
	} else if llm.Config.Temperature != nil {
		body.Temperature = llm.Config.Temperature
	}
	if request.MaxCompletionTokens != nil {
		body.MaxCompletionTokens = request.MaxCompletionTokens
	} else if llm.Config.MaxTokens != nil {
		body.MaxCompletionTokens = llm.Config.MaxTokens
	}
	if request.TopP != nil {
		body.TopP = request.TopP
	} else if llm.Config.TopP != nil {
		body.TopP = llm.Config.TopP
	}
	if request.N != nil {
		body.N = request.N
	} else if llm.Config.N != nil {
		body.N = llm.Config.N
	}
	if request.Stream != nil {
		body.Stream = request.Stream
	} else if llm.Config.Stream != nil {
		body.Stream = llm.Config.Stream
	}

	return body
}

func (llm *LLM) CallCompletion(request types.ChatCompletionRequest) (*types.ChatCompletionResponse, error) {
	headers := http.Header{
		"Content-Type": {"application/json"},
		"User-Agent":   {"SnideMind/LLMClient"},
	}
	if llm.Config.Headers != nil {
		for key, value := range llm.Config.Headers {
			headers.Add(key, value)
		}
	}
	if llm.Config.APIKey != nil && llm.Config.APIKeyHeader != nil {
		headers.Add(*llm.Config.APIKeyHeader, *llm.Config.APIKey)
	}
	urlJoin, err := url.JoinPath(llm.Config.BaseURL, "/v1/chat/completions")
	if err != nil {
		return nil, err
	}
	llmUrl, err := url.Parse(urlJoin)
	if err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(llm.BuildCompletionRequest(request))
	if err != nil {
		return nil, err
	}
	llmRequest := http.Request{
		Method: http.MethodPost,
		URL:    llmUrl,
		Header: headers,
		Body:   io.NopCloser(bytes.NewReader(requestBody)),
	}
	resp, err := http.DefaultClient.Do(&llmRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, LLMError{
			StatusCode: resp.StatusCode,
			Message:    "LLM request failed",
		}
	}
	var response types.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, LLMError{
			StatusCode: resp.StatusCode,
			Message:    "Failed to decode LLM response",
		}
	}
	return &response, nil
}
