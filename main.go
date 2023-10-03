package main

import (
	"fmt"
	"github.com/6691a/iac/config"
	"github.com/6691a/iac/router"
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
	//initSentry(*setting)
	instance, err := router.NewRouter(*setting)
	rtSetting := setting.Router.RouterOS
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(instance.Login(rtSetting.User, rtSetting.Password))

}
