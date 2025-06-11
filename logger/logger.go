package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
}

func NewLogger(p Params) *zap.Logger {
	logCfg := zap.NewProductionConfig()
	logCfg.EncoderConfig.FunctionKey = "method"
	logger := zap.Must(logCfg.Build())

	return logger
}
