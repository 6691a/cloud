package main

import (
	"github.com/6691a/iac/config"
	"github.com/6691a/iac/dns"
	"github.com/getsentry/sentry-go"
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
	initSentry(*setting)
	tasks := dns.GetTasks(*setting)
	dns.CreateWorker(*setting, tasks)
}
