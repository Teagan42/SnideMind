// File: logger/module_test.go
package logger

import (
	"testing"

	"github.com/teagan42/snidemind/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

func TestModule_ProvidesLogger(t *testing.T) {
	app := fxtest.New(
		t,
		Module,
		fx.Invoke(func(l *zap.Logger) {
			if l == nil {
				t.Error("Expected logger to be provided, got nil")
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}

func TestModule_ContainsLoggerModule(t *testing.T) {
	// This test ensures that the Module variable is not nil and is an fx.Option.
	if Module == nil {
		t.Error("Module should not be nil")
	}
}
func TestConfigModule_ProvidesConfig(t *testing.T) {
	app := fxtest.New(
		t,
		fx.Provide(func() *zap.Logger {
			return zap.NewNop()
		}),
		fx.Provide(func() *config.Config {
			return &config.Config{}
		}),
		// Import the config module from the config package.
		fx.Provide(func() fx.Option { return config.Module }),
		fx.Invoke(func(cfg *config.Config) {
			if cfg == nil {
				t.Error("Expected config to be provided, got nil")
			}
		}),
	)
	app.RequireStart()
	app.RequireStop()
}

func TestConfigModule_IsNotNil(t *testing.T) {
	if config.Module == nil {
		t.Error("config.Module should not be nil")
	}
}
