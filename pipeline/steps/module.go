package steps

import (
	extracttags "github.com/teagan42/snidemind/pipeline/steps/extractTags"
	"github.com/teagan42/snidemind/pipeline/steps/fork"
	"github.com/teagan42/snidemind/pipeline/steps/llm"
	reducetools "github.com/teagan42/snidemind/pipeline/steps/reduceTools"
	retrievememory "github.com/teagan42/snidemind/pipeline/steps/retrieveMemory"
	storememory "github.com/teagan42/snidemind/pipeline/steps/storeMemory"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"steps",
	extracttags.Module,
	fork.Module,
	llm.Module,
	reducetools.Module,
	retrievememory.Module,
	storememory.Module,
)
