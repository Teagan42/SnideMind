package chat

import (
	"net/http"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/schema"
	"github.com/teagan42/snidemind/server/middleware"
	"github.com/teagan42/snidemind/server/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ChatCompletionsControllerParams struct {
	fx.In
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
	Config    *config.Config
}

type ChatCompletionsController struct {
	log *zap.Logger
}

func NewChatCompletionsController(p ChatCompletionsControllerParams) *ChatCompletionsController {
	return &ChatCompletionsController{
		log: p.Log,
	}
}

func (c *ChatCompletionsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, err := middleware.GetValidatedBody[*schema.ChatCompletionRequest](r); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// TODO
}

func (c *ChatCompletionsController) Pattern() string {
	return "completions"
}

func (c *ChatCompletionsController) Methods() []string {
	return []string{http.MethodPost}
}

var _ utils.Route = (*ChatCompletionsController)(nil)
