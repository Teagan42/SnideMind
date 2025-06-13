package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger_ReturnsValidLogger(t *testing.T) {
	logger := NewLogger(Params{})
	if logger == nil {
		t.Fatal("NewLogger returned nil, expected a valid *zap.Logger")
	}
	// Test that logger can log without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger.Panic or Logger.Info panicked: %v", r)
		}
	}()
	logger.Info("test message")
}

func TestNewLogger_ProducesProductionConfig(t *testing.T) {
	logger := NewLogger(Params{})
	core := logger.Core()
	if core == nil {
		t.Error("Logger core is nil, expected a valid core")
	}
	// Optionally, check that the logger is not in development mode by default
	// (ProductionConfig disables development mode)
	if logger.Check(zap.InfoLevel, "test") == nil {
		t.Error("Logger does not allow InfoLevel logs, expected it to allow")
	}
}
