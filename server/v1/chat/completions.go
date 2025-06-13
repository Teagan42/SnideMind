package chat

import (
	"encoding/json"
	"net/http"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/pipeline"
	"github.com/teagan42/snidemind/server/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ChatCompletionsControllerParams struct {
	fx.In
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
	Config    *config.Config
	Pipeline  *pipeline.Pipeline
}

type ChatCompletionsController struct {
	log      *zap.Logger
	pipeline *pipeline.Pipeline
}

func NewChatCompletionsController(p ChatCompletionsControllerParams) *ChatCompletionsController {
	return &ChatCompletionsController{
		log:      p.Log.Named("ChatCompletionsController"),
		pipeline: p.Pipeline,
	}
}

func (c *ChatCompletionsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	c.log.Info("Processing pipeline")
	message, err := c.pipeline.Process(r)
	if err != nil {
		c.log.Error("Error processing pipeline", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	c.log.Info("Pipeline processed successfully", zap.Any("message", message))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody, err := json.Marshal(message)
	if err != nil {
		c.log.Error("Error unmarshalling response", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(respBody); err != nil {
		c.log.Error("Error writing JSON response", zap.Error(err))
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func (c *ChatCompletionsController) Pattern() string {
	return "completions"
}

func (c *ChatCompletionsController) Methods() []string {
	return []string{http.MethodPost}
}

var _ utils.Route = (*ChatCompletionsController)(nil)
