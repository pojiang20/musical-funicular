package util

import "go.uber.org/zap"

var Zlog *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	Zlog = logger.Sugar()
}
