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

	instance.Delete(1050)
	instance.Delete(1051)
	cloneRecode := hypervisor.NewCloneRecord(
		10002,
		"golang.test.com",
		"Test proxmox-api-go clone",
		1050,
	)
	// Clone VM
	fmt.Print(instance.Clone(cloneRecode))

	vcrd := hypervisor.NewVmConfigRecode(
		1050,
		1,
		1024,
		true,
		map[uint8]string{0: "virtio,bridge=vmbr2"},
		map[uint8]string{0: "gw=192.168.10.1,ip=192.168.10.2/24"},
	)
	// Set VM config
	fmt.Print(instance.SetVmConfig(vcrd))

	// ========================
	cloneRecode = hypervisor.NewCloneRecord(
		10002,
		"golang.test.com",
		"Test proxmox-api-go clone",
		1051,
	)
	// Clone VM
	fmt.Print(instance.Clone(cloneRecode))

	vcrd = hypervisor.NewVmConfigRecode(
		1051,
		1,
		1024,
		true,
		map[uint8]string{0: "virtio,bridge=vmbr2"},
		map[uint8]string{0: "gw=192.168.10.1,ip=192.168.10.10/24"},
	)
	// Set VM config
	fmt.Print(instance.SetVmConfig(vcrd))

}
