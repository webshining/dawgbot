package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}
