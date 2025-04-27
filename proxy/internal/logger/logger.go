package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()
}
