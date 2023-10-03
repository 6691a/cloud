package main

import (
	"fmt"
	"github.com/6691a/iac/config"
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
	tasks := hypervisor.GetTasks(*setting)
	hypervisor.CreateWorker(*setting, tasks)

	cloneRecode := hypervisor.NewCloneRecord(
		10002,
		"golang.test.com",
		"Test proxmox-api-go clone",
		1050,
	)
	cloneTask := hypervisor.Task{
		Request: hypervisor.Request{
			Method: hypervisor.Clone,
			Record: cloneRecode,
		},
		Response: make(chan hypervisor.Response),
	}
	tasks <- cloneTask

	// 응답 대기
	cloneResponse := <-cloneTask.Response
	fmt.Print(cloneResponse.Error)

	vcrd := hypervisor.NewVmConfigRecode(
		1050,
		1,
		1024,
		false,
		map[uint8]string{
			0: "virtio,bridge=vmbr0",
		},
		map[uint8]string{
			0: "ip=dhcp",
		},
	)
	setVmTask := hypervisor.Task{
		Request: hypervisor.Request{
			Method: hypervisor.SetVmConfig,
			Record: vcrd,
		},
		Response: make(chan hypervisor.Response),
	}
	tasks <- setVmTask

	// 응답 대기
	setVmResponse := <-setVmTask.Response
	fmt.Print(setVmResponse.Error)
}
