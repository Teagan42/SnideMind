package models

import (
	"net/http"

	"github.com/teagan42/snidemind/server/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ListModelsController struct {
	log *zap.Logger
}

type ListModelsControllerParams struct {
	fx.In
	Log       *zap.Logger
	Lifecycle fx.Lifecycle
}

func NewListModelsController(p ListModelsControllerParams) *ListModelsController {
	return &ListModelsController{
		log: p.Log,
	}
}

func (c *ListModelsController) Pattern() string {
	return "/"
}

func (c *ListModelsController) Methods() []string {
	return []string{http.MethodGet}
}

func (c *ListModelsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Placeholder for the actual implementation
	// This should list available models and return them in the response
	c.log.Info("ListModelsController ServeHTTP called", zap.String("method", r.Method), zap.String("url", r.URL.String()))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of models will be implemented here"))
}

var _ utils.Route = (*ListModelsController)(nil)
