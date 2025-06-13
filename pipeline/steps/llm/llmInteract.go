package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/schema"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LLM struct {
	config.LLMConfig
	Logger *zap.Logger
}

type Params struct {
	fx.In
	Logger *zap.Logger
}

type Result struct {
	fx.Out
	Factory schema.PipelineStepFactory `group:"pipelineStepFactory"`
}

type LLMFactory struct {
	Logger *zap.Logger
}

func (f LLMFactory) Name() string {
	return "llm"
}
func (f LLMFactory) Build(config config.PipelineStepConfig) (schema.PipelineStep, error) {
	return &LLM{
		LLMConfig: *config.LLM,
		Logger:    f.Logger,
	}, nil
}

func NewLLM(p Params) (Result, error) {
	return Result{
		Factory: LLMFactory{
			Logger: p.Logger.Named("LLMInteract"),
		},
	}, nil
}

func (s LLM) Name() string {
	return "LLM"
}

func (s LLM) Process(previous *[]schema.PipelineStep, input *schema.PipelineMessage) (*schema.PipelineMessage, error) {
	s.Logger.Info("Processing", zap.String("model", *s.Model), zap.String("baseURL", s.BaseURL))
	bodyBytes, err := json.Marshal(input.Request)
	if err != nil {
		s.Logger.Error("Error marshalling request", zap.Error(err))
		return nil, err
	}
	url := fmt.Sprintf("%s/chat/completions", s.BaseURL)
	s.Logger.Info("Creating request", zap.String("url", url), zap.ByteString("body", bodyBytes))
	client := &http.Client{}
	if req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewBuffer(bodyBytes))); err != nil {
		s.Logger.Error("Error creating request", zap.Error(err))
		return nil, err
	} else {
		if s.APIKey != nil && s.APIKeyHeader != nil {
			req.Header.Set(*s.APIKeyHeader, *s.APIKey)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		s.Logger.Info("Sending request", zap.String("url", url), zap.ByteString("body", bodyBytes))
		if resp, err := client.Do(req); err != nil {
			s.Logger.Error("Error sending request", zap.Error(err))
			return nil, err
		} else {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				s.Logger.Error("Error response from LLM", zap.String("status", resp.Status))
				return nil, fmt.Errorf("error: %s", resp.Status)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				s.Logger.Error("Error reading response body", zap.Error(err))
				return nil, err
			}
			var response schema.ChatCompletionResponse
			if err := json.Unmarshal(body, &response); err != nil {
				s.Logger.Error("Error unmarshalling response", zap.Error(err))
				return nil, err
			}
			s.Logger.Info("Received response", zap.String("response", string(body)))
			input.Response = &response
		}
	}

	return input, nil
}
