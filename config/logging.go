package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = make(map[string]*zap.Logger)

func NewLogging(setting Setting) {
	serverSetting := setting.Server

	var encodeConfig zapcore.EncoderConfig
	var config zap.Config

	if serverSetting.Debug {
		encodeConfig = zap.NewDevelopmentEncoderConfig()
		config = zap.NewDevelopmentConfig()
	} else {
		encodeConfig = zap.NewProductionEncoderConfig()
		config = zap.NewProductionConfig()

	}

	for name, loggingSetting := range serverSetting.Logging {
		encodeConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encodeConfig.TimeKey = "timestamp"
		encodeConfig.LevelKey = "level"
		encodeConfig.MessageKey = "event"
		encodeConfig.CallerKey = "pathname"
		encodeConfig.FunctionKey = "func_name"
		encodeConfig.LineEnding = zapcore.DefaultLineEnding
		encodeConfig.EncodeDuration = zapcore.StringDurationEncoder
		encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encodeConfig.EncodeName = zapcore.FullNameEncoder

		config.EncoderConfig = encodeConfig
		config.Encoding = "json"
		config.OutputPaths = loggingSetting.Path

		logger[name] = zap.Must(config.Build())
	}
}

func GetLogger(name string) *zap.Logger {

	if logger == nil {
		panic("Logger is not initialized")
	}

	if _, ok := logger[name]; !ok {
		logger[name] = zap.L().Named(name)
	}

	return logger[name]
}
