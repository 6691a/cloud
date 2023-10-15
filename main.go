package main

import (
	"github.com/6691a/iac/config"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func initSentry(setting config.Setting) {
	if !setting.Server.Debug {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              setting.Server.SentryDsn,
			TracesSampleRate: 1.0,
		})

		if err != nil {
			panic(err)
		}
	}
}

func main() {
	setting := config.NewSetting("setting.yaml")
	config.NewLogging(*setting)
	logger := config.GetLogger("default")
	logger.Info("user registration successful",
		zap.String("username", "john.doe"),
		zap.String("email", "john@example.com"),
	)

}
