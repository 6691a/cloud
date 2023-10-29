package config

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var logger = make(map[string]*zap.Logger)

func NewLogging(setting Setting) {
	serverSetting := setting.Server
	for name, loggingSetting := range serverSetting.Logging {
		config, err := createLoggerConfig(serverSetting.Debug, loggingSetting)
		if err != nil {
			panic(err)
		}
		logger[name] = createLogger(config, serverSetting.Debug, loggingSetting)
	}
}

func createLoggerConfig(isDebug bool, loggingSetting LoggingConfig) (zap.Config, error) {
	var config zap.Config
	encodeConfig := createEncoderConfig()

	if isDebug {
		config = zap.NewDevelopmentConfig()
		config.Encoding = "console"
	} else {
		config = zap.NewProductionConfig()
		config.Encoding = "json"
	}

	config.EncoderConfig = encodeConfig
	config.OutputPaths = []string{loggingSetting.Path}

	return config, nil
}

func createEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "event",
		CallerKey:      "pathname",
		FunctionKey:    "func_name",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func createLogger(config zap.Config, isDebug bool, loggingSetting LoggingConfig) *zap.Logger {
	var core zapcore.Core
	encodeConfig := config.EncoderConfig

	if isDebug {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encodeConfig),
			zapcore.AddSync(zapcore.Lock(os.Stdout)),
			zapcore.DebugLevel,
		)
	} else {
		logRotate := createFileRotate(loggingSetting) // 이 줄에 오류가 있습니다. loggingSetting은 이 스코프에서 사용할 수 없습니다.
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encodeConfig),
			zapcore.AddSync(logRotate),
			zapcore.DebugLevel,
		)
	}

	return zap.New(core, zap.AddCaller())
}

func createFileRotate(config LoggingConfig) *rotatelogs.RotateLogs {
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
