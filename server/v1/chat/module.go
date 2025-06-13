package chat

import (
	"github.com/teagan42/snidemind/server/utils"
)

var Module = utils.ApiRouteModule(utils.ApiRouteModuleParams{
	ParentRouter: "v1",
	ModuleName:   "chat",
	Prefix:       "chat",
	Routes: &[]any{
		NewChatCompletionsController,
	},
})
