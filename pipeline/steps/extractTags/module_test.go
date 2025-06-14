package extracttags

import (
	"testing"

	"github.com/teagan42/snidemind/models"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestExtractTagsModule_IsNotNil(t *testing.T) {
	// This test ensures that the extracttags.Module variable is not nil and is an fx.Option.
	if extractTagsModule := fx.Option(Module); extractTagsModule == nil {
		t.Error("extracttags.Module should not be nil")
	}
}

type TestParams struct {
	fx.In
	Factory []models.PipelineStepFactory `group:"pipelineStepFactory"`
}

func TestExtractTagsModule_ProvidesExtractTags(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		Module,
		fx.Invoke(func(p TestParams) {
			if p.Factory == nil {
				t.Error("Expected ExtractTags to be provided, got nil")
			}
			if len(p.Factory) == 0 {
				t.Error("Expected ExtractTags to be provided, got empty slice")
			}
			if p.Factory[0].Name() != "extractTags" {
				t.Errorf("Expected ExtractTags factory name to be 'extractTags', got '%s'", p.Factory[0].Name())
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}
