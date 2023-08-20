package main

import (
	"fmt"
	"github.com/6691a/iac/config"
	"github.com/6691a/iac/dns"
	"github.com/6691a/iac/hypervisor"
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

	instance, err := hypervisor.NewHypervisor(*setting)
	if err != nil {
		fmt.Print(err)
	}
	cloneRecode := hypervisor.NewCloneRecord(10002, "golang.test.com", "Test proxmox-api-go clone", 1051)

	instance.Clone(cloneRecode)
	instance.Delete(1051)

}
