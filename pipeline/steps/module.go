package steps

import (
	extracttags "github.com/teagan42/snidemind/pipeline/steps/extractTags"
	"github.com/teagan42/snidemind/pipeline/steps/fork"
	"github.com/teagan42/snidemind/pipeline/steps/llm"
	reducetools "github.com/teagan42/snidemind/pipeline/steps/reduceTools"
	retrievememory "github.com/teagan42/snidemind/pipeline/steps/retrieveMemory"
	storememory "github.com/teagan42/snidemind/pipeline/steps/storeMemory"
	"github.com/teagan42/snidemind/schema"
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
	fx.Provide(
		fx.Annotate(
			func(stepFactories []schema.PipelineStepFactory) map[string]schema.PipelineStepFactory {
				stepMap := make(map[string]schema.PipelineStepFactory)
				for _, factory := range stepFactories {
					stepMap[factory.Name()] = factory
				}
				return stepMap
			},
			fx.ParamTags(`group:"pipelineStepFactory"`),
			fx.ResultTags(`name:"pipelineStepFactoryMap"`),
		),
	),
)
