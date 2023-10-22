package config

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
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
		logRotate := createFileRotate(loggingSetting)

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
		config.OutputPaths = []string{loggingSetting.Path}

		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(encodeConfig),
				zapcore.AddSync(logRotate),
				zapcore.DebugLevel,
			),
		)
		logger[name] = zap.New(core)
	}
}

func createFileRotate(config LoggingConfig) *rotatelogs.RotateLogs {
	fmt.Println(config.Path)
	fileRotate, err := rotatelogs.New(
		config.Path+".%Y%m%d.log",
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)*time.Hour),
		rotatelogs.WithMaxAge(time.Duration(config.MaxAge)*24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return fileRotate
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
