package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"github.com/teagan42/snidemind/pipeline"
	"github.com/teagan42/snidemind/server/middleware"
	"go.uber.org/zap"
)

type MockStep struct {
	called bool
}

func (m *MockStep) Name() string {
	return "MockStep"
}
func (m *MockStep) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	fmt.Printf("MockStep called with input: %v\n", input)
	input.ResponseWriter.Write([]byte("MockStep called"))
	m.called = true
	input.ResponseWriter.WriteHeader(http.StatusOK)
	input.ResponseWriter.Header().Set("Content-Type", "application/json")
	fmt.Println("MockStep done")
	return input, nil
}

func NewMockStep() *MockStep {
	return &MockStep{}
}

var mockStep = &MockStep{}

type MockStepFactory struct{}

func (f MockStepFactory) Name() string {
	return "MockStep"
}

func (f MockStepFactory) Build(config config.PipelineStepConfig, stepFactories map[string]models.PipelineStepFactory) (models.PipelineStep, error) {
	return mockStep, nil
}

var testConfig = &config.Config{
	Server: config.ServerConfig{
		Port: 8080,
	},
	Pipeline: &config.PipelineConfig{
		Steps: []config.PipelineStepConfig{
			{
				Type: "MockStep",
			},
		},
	},
}

func TestChatCompletionsController_ServeHTTP_MethodNotAllowed(t *testing.T) {
	mockStep.called = false
	ctrl := NewChatCompletionsController(ChatCompletionsControllerParams{
		Log:    zap.L().Named("TestChatCompletionsController"),
		Config: testConfig,
		Pipeline: pipeline.NewPipeline(pipeline.Params{
			Config: testConfig,
			Logger: zap.L().Named("TestChatCompletionsController"),
			StepFactories: map[string]models.PipelineStepFactory{
				"MockStep": MockStepFactory{},
			},
		}),
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	ctrl.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestChatCompletionsController_ServeHTTP_PipelineError(t *testing.T) {
	mockStep.called = false
	ctrl := &ChatCompletionsController{
		log: zap.NewNop(),
		pipeline: pipeline.NewPipeline(pipeline.Params{
			Config: testConfig,
			Logger: zap.NewNop(),
		}),
	}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{}`))
	rr := httptest.NewRecorder()

	ctrl.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
	want := `{"error":"internal server error"}`
	if rr.Body.String() != want+"\n" && rr.Body.String() != want {
		t.Errorf("expected body %q, got %q", want, rr.Body.String())
	}
}

func TestChatCompletionsController_ServeHTTP_PipelineSuccess(t *testing.T) {
	mockStep.called = false
	ctrl := &ChatCompletionsController{
		log: zap.NewNop(),
		pipeline: pipeline.NewPipeline(pipeline.Params{
			Config: testConfig,
			Logger: zap.NewNop(),
			StepFactories: map[string]models.PipelineStepFactory{
				"MockStep": MockStepFactory{},
			},
		}),
	}
	body := models.ChatCompletionRequest{
		Messages: []models.ChatMessage{
			{
				Role:    "user",
				Content: "Hello, world!",
			},
		},
		Model: "gpt-3.5-turbo",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("failed to marshal request body: %v", err)
		return
	}
	var raw any
	if err := json.NewDecoder(io.NopCloser(bytes.NewBuffer(bodyBytes))).Decode(&raw); err != nil {
		t.Errorf("failed to decode request body: %v", err)
		return
	}
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBuffer(bodyBytes)))
	req = req.WithContext(
		context.WithValue(
			req.Context(),
			middleware.BodyKey,
			raw,
		),
	)
	rr := httptest.NewRecorder()
	ctrl.ServeHTTP(rr, req)

	time.Sleep(1 * time.Second) // Wait for the goroutine to finish

	if !mockStep.called {
		t.Error("expected pipeline.Process to be called")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
