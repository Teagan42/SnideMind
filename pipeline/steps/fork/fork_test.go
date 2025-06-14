// File: pipeline/steps/fork/fork_test.go
package fork

import (
	"errors"
	"testing"

	"github.com/teagan42/snidemind/config"
	"github.com/teagan42/snidemind/models"
	"go.uber.org/zap/zaptest"
)

type dummyStep struct {
	models.PipelineStep
	name string
}

func (d *dummyStep) Name() string { return d.name }
func (d *dummyStep) Process(_ *[]models.PipelineStep, msg *models.PipelineMessage) (*models.PipelineMessage, error) {
	return msg, nil
}

type dummyFactory struct {
	models.PipelineStepFactory
	buildErr error
	stepType string
}

func (d dummyFactory) Name() string { return d.stepType }
func (d dummyFactory) Build(_ config.PipelineStepConfig, _ map[string]models.PipelineStepFactory) (models.PipelineStep, error) {
	if d.buildErr != nil {
		return nil, d.buildErr
	}
	return &dummyStep{name: d.stepType}, nil
}

func TestForkPipelineStageFactory_Build_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := ForkPipelineStageFactory{Logger: logger}

	stepType := "dummy"
	stepConfig := config.PipelineStepConfig{Type: stepType}
	forkConfig := []config.PipelineConfig{
		{
			Steps: []config.PipelineStepConfig{stepConfig},
		},
	}
	cfg := config.PipelineStepConfig{
		Fork: &forkConfig,
	}

	stepFactories := map[string]models.PipelineStepFactory{
		stepType: dummyFactory{stepType: stepType},
	}

	stage, err := factory.Build(cfg, stepFactories)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stage == nil {
		t.Fatal("expected non-nil stage")
	}
	forkStage, ok := stage.(*ForkPipelineStage)
	if !ok {
		t.Fatalf("expected *ForkPipelineStage, got %T", stage)
	}
	if len(forkStage.Forks) != 1 {
		t.Errorf("expected 1 fork, got %d", len(forkStage.Forks))
	}
	if len(forkStage.Forks[0].Steps) != 1 {
		t.Errorf("expected 1 step in fork, got %d", len(forkStage.Forks[0].Steps))
	}
}

func TestForkPipelineStageFactory_Build_NilForkConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := ForkPipelineStageFactory{Logger: logger}
	cfg := config.PipelineStepConfig{Fork: nil}
	_, err := factory.Build(cfg, nil)
	if err == nil || err.Error() != "fork config is nil" {
		t.Errorf("expected error 'fork config is nil', got %v", err)
	}
}

func TestForkPipelineStageFactory_Build_EmptyForkConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := ForkPipelineStageFactory{Logger: logger}
	empty := []config.PipelineConfig{}
	cfg := config.PipelineStepConfig{Fork: &empty}
	_, err := factory.Build(cfg, nil)
	if err == nil || err.Error() != "fork config is empty" {
		t.Errorf("expected error 'fork config is empty', got %v", err)
	}
}

func TestForkPipelineStageFactory_Build_EmptyStepType(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := ForkPipelineStageFactory{Logger: logger}
	stepConfig := config.PipelineStepConfig{Type: ""}
	forkConfig := []config.PipelineConfig{
		{
			Steps: []config.PipelineStepConfig{stepConfig},
		},
	}
	cfg := config.PipelineStepConfig{Fork: &forkConfig}
	stepFactories := map[string]models.PipelineStepFactory{}
	_, err := factory.Build(cfg, stepFactories)
	if err == nil || err.Error() != "step type is empty in fork config at index 0" {
		t.Errorf("expected error for empty step type, got %v", err)
	}
}

func TestForkPipelineStageFactory_Build_UnknownStepType(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := ForkPipelineStageFactory{Logger: logger}
	stepConfig := config.PipelineStepConfig{Type: "unknown"}
	forkConfig := []config.PipelineConfig{
		{
			Steps: []config.PipelineStepConfig{stepConfig},
		},
	}
	cfg := config.PipelineStepConfig{Fork: &forkConfig}
	stepFactories := map[string]models.PipelineStepFactory{}
	_, err := factory.Build(cfg, stepFactories)
	if err == nil || err.Error() != "unknown pipeline step type: unknown" {
		t.Errorf("expected error for unknown step type, got %v", err)
	}
}

func TestForkPipelineStageFactory_Build_FactoryBuildError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	factory := ForkPipelineStageFactory{Logger: logger}
	stepType := "dummy"
	stepConfig := config.PipelineStepConfig{Type: stepType}
	forkConfig := []config.PipelineConfig{
		{
			Steps: []config.PipelineStepConfig{stepConfig},
		},
	}
	cfg := config.PipelineStepConfig{Fork: &forkConfig}
	stepFactories := map[string]models.PipelineStepFactory{
		stepType: dummyFactory{stepType: stepType, buildErr: errors.New("fail")},
	}
	_, err := factory.Build(cfg, stepFactories)
	if err == nil || err.Error() != "failed to build pipeline step: fail" {
		t.Errorf("expected error for failed build, got %v", err)
	}
}
