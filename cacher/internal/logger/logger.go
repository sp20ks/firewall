package logger

import (
	"context"
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

var logger *zap.Logger

type ctxLogger struct{}

func init() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger = zap.New(core, zap.AddCaller()).With(zap.String("service", "cacher")).With(zap.String("environment", "development"))
}

func Logger() *zap.Logger {
	return logger
}

// do not use.
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}
