package reducetools

import (
	"maps"
	"slices"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"github.com/teagan42/snidemind/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ReduceTools struct {
	Logger  *zap.Logger
	ToolSet []models.MCPTool
}

type Params struct {
	fx.In
	Logger *zap.Logger
}

type Result struct {
	fx.Out
	Factory models.PipelineStepFactory `group:"pipelineStepFactory"`
}

type ReduceToolsFactory struct {
	Logger  *zap.Logger
	ToolSet []models.MCPTool
}

func (f ReduceToolsFactory) Name() string {
	return "reduceTools"
}
func (f ReduceToolsFactory) Build(config config.PipelineStepConfig, stepFactories map[string]models.PipelineStepFactory) (models.PipelineStep, error) {
	return &ReduceTools{
		Logger:  f.Logger.Named("ReduceTools"),
		ToolSet: f.ToolSet,
	}, nil
}

func NewReduceTools(p Params) (Result, error) {
	return Result{
		Factory: ReduceToolsFactory{
			Logger: p.Logger.Named("ReduceToolsFactory"),
			ToolSet: []models.MCPTool{
				{
					ToolMetadata: models.ToolMetadata{
						Name:        "Plex New Media",
						Description: "Retrieves new media from Plex",
						Tags:        &[]string{"plex", "media", "media.new"},
					},
				},
				{
					ToolMetadata: models.ToolMetadata{
						Name:        "Plex Search",
						Description: "Searches Plex for media",
						Tags:        &[]string{"plex", "media", "media.search", "media.movies", "media.tv"},
					},
				},
				{
					ToolMetadata: models.ToolMetadata{
						Name:        "Home Assistant Entity Action",
						Description: "Performs an action on a Home Assistant entity",
						Tags:        &[]string{"home", "home.automation"},
					},
				},
			},
		},
	}, nil
}

func (s ReduceTools) Name() string {
	return "reduceTools"
}

func (s ReduceTools) Process(previous *[]models.PipelineStep, input *models.PipelineMessage) (*models.PipelineMessage, error) {
	if input.Tags == nil || len(*input.Tags) == 0 {
		s.Logger.Debug("No tags found, skipping tool reduction")
		return input, nil
	}
	tags := slices.Collect(maps.Values(*input.Tags))
	input.Tools = &[]models.MCPTool{}
	for _, tool := range s.ToolSet {
		s.Logger.Debug("Checking tool", zap.String("tool", tool.ToolMetadata.Name))

		if len(utils.Intersection(tags, *tool.ToolMetadata.Tags)) > 0 {
			s.Logger.Debug("Tool matches tags", zap.String("tool", tool.ToolMetadata.Name))
			*input.Tools = append(*input.Tools, tool)
		} else {
			s.Logger.Debug("Tool does not match tags", zap.String("tool", tool.ToolMetadata.Name))
		}
	}

	return input, nil
}
