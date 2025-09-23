package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	// logger, _ := zap.NewDevelopment()
	// zap.ReplaceGlobals(logger)

	cfg := zap.NewDevelopmentConfig()

	// 时间格式：2025-09-21 17:40:01
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}
