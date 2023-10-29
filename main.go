package main

import (
	"fmt"
	"github.com/6691a/iac/config"
	"github.com/6691a/iac/dns"
	"github.com/6691a/iac/hypervisor"
	"github.com/6691a/iac/router"
	"github.com/6691a/iac/router/routeros"
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

func setup() {
	setting := config.NewSetting("setting.yaml")
	config.NewLogging(*setting)
	initSentry(*setting)
}

func main() {
	setup()
	setting := config.GetSetting()

	dnsTasks := dns.GetTasks(*setting)
	dns.CreateWorker(*setting, dnsTasks)
	hvTasks := hypervisor.GetTasks(*setting)
	hypervisor.CreateWorker(*setting, hvTasks)
	router, _ := router.NewRouter(*setting)

	// ==============================

	// ============ CLONE ============
	//Clone(tasks)

	// ============ CONFIG ============
	//Config(hvTasks)

	// ============ DELETE ============
	//Delete(tasks)

	// ============ DNS ============
	//CreateDNS(dnsTasks)

	// ============ ADD WITH LIST ============
	//CreateWhiteList(router)
	GetWithList(router)
	//CreateDomainList(router)

}

func Clone(task chan hypervisor.Task) {
	cloneRecode := hypervisor.NewCloneRecord(
		10002,
		"sample",
		"sample clone",
		1050,
	)
	cloneTask := hypervisor.Task{
		Request: hypervisor.Request{
			Method: hypervisor.Clone,
			Record: cloneRecode,
		},
		Response: make(chan hypervisor.Response),
	}
	task <- cloneTask

	//응답 대기
	cloneResponse := <-cloneTask.Response
	fmt.Print(cloneResponse.Error)
}

func Delete(task chan hypervisor.Task) {
	deleteTask := hypervisor.Task{
		Request: hypervisor.Request{
			Method: hypervisor.Delete,
			Record: hypervisor.NewRecord(1050, "", ""),
		},
		Response: make(chan hypervisor.Response),
	}

	task <- deleteTask
	fmt.Println(<-deleteTask.Response)
}

func Config(task chan hypervisor.Task) {
	vcrd := hypervisor.NewVmConfigRecode(
		1050,
		1,
		1024,
		false,
		map[uint8]string{
			0: "virtio,bridge=vmbr0",
		},
		map[uint8]string{
			0: "gw=192.168.88.1,ip=192.168.88.254/24",
		},
	)
	setVmTask := hypervisor.Task{
		Request: hypervisor.Request{
			Method: hypervisor.SetVmConfig,
			Record: vcrd,
		},
		Response: make(chan hypervisor.Response),
	}
	task <- setVmTask

	// 응답 대기
	setVmResponse := <-setVmTask.Response
	fmt.Print(setVmResponse)
}

func CreateDNS(task chan dns.Task) {
	createTask := dns.Task{
		Request: dns.Request{
			Method: dns.Method(dns.Create),
			Record: *dns.NewRecord("sample", "CNAME", 300, []string{"playhub.kr."}),
		},
		Response: make(chan dns.Response),
	}
	task <- createTask

	createResponse := <-createTask.Response
	fmt.Printf("Received create response: %+v\n", createResponse)
}

func CreateWhiteList(router router.Router) {
	rtOS := router.(*routeros.RouterOS)
	fmt.Println(rtOS.CreateWhiteList("sample", "192.168.100.100"))
}

func GetWithList(router router.Router) {
	rtOS := router.(*routeros.RouterOS)
	fmt.Println(rtOS.GetWhiteList("sample1", "192.168.100.100"))
}

func CreateDomainList(router router.Router) {
	rtOS := router.(*routeros.RouterOS)
	fmt.Println(rtOS.CreateDomainList("sample", "playhub.kr"))
}
