package v1

import (
	"github.com/teagan42/snidemind/server/utils"
	"github.com/teagan42/snidemind/server/v1/chat"
	"github.com/teagan42/snidemind/server/v1/models"
	"go.uber.org/fx"
)

var Module = utils.ApiRouteModule(utils.ApiRouteModuleParams{
	ParentRouter: "root",
	ModuleName:   "v1",
	Prefix:       "v1",
	SubModules: &[]fx.Option{
		chat.Module,
		models.Module,
	},
})
