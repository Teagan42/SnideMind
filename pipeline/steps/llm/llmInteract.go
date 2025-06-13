package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
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
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type LLMFactory struct {
	Logger *zap.Logger
}

func (f LLMFactory) Name() string {
	return "llm"
}
func (f LLMFactory) Build(config config.PipelineStepConfig, stepFactories map[string]models.PipelineStepFactory) (models.PipelineStep, error) {
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

func (s LLM) buildRequestBody(input *models.PipelineMessage) models.ChatCompletionRequest {
	reqBody := models.ChatCompletionRequest{
		Messages:         input.Request.Messages,
		Model:            input.Request.Model,
		FrequencyPenalty: input.Request.FrequencyPenalty,
		N:                input.Request.N,
		// Tools:             input.Tools,
		ParallelToolCalls: input.Request.ParallelToolCalls,
		PresencePenalty:   input.Request.PresencePenalty,
		Stream:            input.Request.Stream,
		Temperature:       input.Request.Temperature,
		TopP:              input.Request.TopP,
	}
	if s.FrequencyPenalty != nil {
		reqBody.FrequencyPenalty = s.FrequencyPenalty
	}
	if s.N != nil {
		reqBody.N = s.N
	}
	if s.ParallelToolCalls != nil {
		reqBody.ParallelToolCalls = s.ParallelToolCalls
	}
	if s.PresencePenalty != nil {
		reqBody.PresencePenalty = s.PresencePenalty
	}
	if s.Stream != nil {
		reqBody.Stream = s.Stream
	}
	if s.Temperature != nil {
		reqBody.Temperature = s.Temperature
	}
	if s.TopP != nil {
		reqBody.TopP = s.TopP
	}

	if input.Tools != nil && len(*input.Tools) > 0 {
		tools := make([]models.Tool, len(*input.Tools))
		for i, tool := range *input.Tools {
			tools[i] = models.Tool{
				Type: "function",
				Function: models.ToolFunction{
					Name:        tool.ToolMetadata.Name,
					Description: tool.ToolMetadata.Description,
					Parameters:  models.ToolFunctionParameters{},
				},
			}
		}
		reqBody.Tools = &tools
	}

	return reqBody
}

func (s LLM) streamResponse(input *models.PipelineMessage, body io.Reader) (*models.PipelineMessage, error) {
	w := input.ResponseWriter
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	var resp *models.ChatCompletionResponse = nil
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		s.Logger.Info("Stream chunk", zap.String("chunk", line))
		if resp == nil {
			resp = &models.ChatCompletionResponse{
				ID: "stream",
				Choices: []models.ChatCompletionChoice{
					{
						FinishReason: "stop",
						Index:        0,
						Message: models.ChatMessage{
							Role:    "assistant",
							Content: "",
						},
					},
				},
				Created: 0,
				Model:   *s.Model,
				Object:  "chat.completion",
			}
		}
		if line != "data: [DONE]" {
			if len(line) > 6 && line[:6] == "data: " {
				resp.Choices[0].Message.Content = resp.Choices[0].Message.Content + line[6:]
			}
		}
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			s.Logger.Error("Write error during stream", zap.Error(err))
			return nil, err
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	if err := scanner.Err(); err != nil {
		s.Logger.Error("Stream scanner error", zap.Error(err))
		return nil, err
	}

	input.Response = resp

	w.WriteHeader(http.StatusOK)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	return input, nil
}

func (s LLM) bufferResponse(input *models.PipelineMessage, body io.Reader) (*models.PipelineMessage, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		s.Logger.Error("Read error", zap.Error(err))
		return nil, err
	}
	s.Logger.Info("Buffered response", zap.ByteString("data", data))

	w := input.ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		s.Logger.Error("Write error", zap.Error(err))
		return nil, err
	}

	var resp models.ChatCompletionResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		s.Logger.Error("Unmarshal error", zap.Error(err))
		return nil, err
	}
	input.Response = &resp

	return input, nil
}

func (s LLM) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	s.Logger.Info("Processing", zap.String("model", *s.Model), zap.String("baseURL", s.BaseURL))

	bodyBytes, err := json.Marshal(s.buildRequestBody(input))
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
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
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
			var respMsg *models.PipelineMessage
			if (s.Stream != nil && *s.Stream) || (input.Request.Stream != nil && *input.Request.Stream) {
				respMsg, err = s.streamResponse(input, resp.Body)
			} else {
				respMsg, err = s.bufferResponse(input, resp.Body)
			}
			if err != nil {
				s.Logger.Error("Error reading response body", zap.Error(err))
				return nil, err
			}
			return respMsg, nil
		}
	}
}
