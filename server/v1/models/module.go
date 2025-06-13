package models

import (
	"github.com/teagan42/snidemind/server/utils"
)

var Module = utils.ApiRouteModule(utils.ApiRouteModuleParams{
	ParentRouter: "v1",
	ModuleName:   "models",
	Prefix:       "models",
	Routes: &[]any{
		NewListModelsController,
		NewGetModelController,
	},
})
