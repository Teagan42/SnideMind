package models

import (
	"net/http"

	"github.com/teagan42/snidemind/server/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type GetModelController struct {
	log *zap.Logger
}

type GetModelControllerParams struct {
	fx.In
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
}

type GetModelControllerResult struct {
	fx.Out
	Controller *GetModelController
}

func NewGetModelController(p GetModelControllerParams) *GetModelController {
	return &GetModelController{
		log: p.Log,
	}
}

func (c *GetModelController) Pattern() string {
	return "/{modelId}"
}

func (c *GetModelController) Methods() []string {
	return []string{http.MethodGet}
}

func (c *GetModelController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Placeholder for the actual implementation
	// This should list available models and return them in the response
	c.log.Info("ListModelsController ServeHTTP called", zap.String("method", r.Method), zap.String("url", r.URL.String()))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of models will be implemented here"))
}

var _ utils.Route = (*GetModelController)(nil)
