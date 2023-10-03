package dns

import (
	"errors"
	"fmt"
	"github.com/6691a/iac/config"
	"sync"
)

type Method string

const (
	Get    Method = "GET"
	List   Method = "LIST"
	Create Method = "CREATE"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

type Request struct {
	Method Method
	Record Record
}

type Response struct {
	Error   error
	Records []Record
}

type Task struct {
	Request  Request
	Response chan Response
}

func CreateWorker(setting config.Setting, tasks chan Task) {
	dns, err := NewDNS(setting)
	if err != nil {
		panic(err)
	}

	dnsSetting := setting.DNS
	size := dnsSetting.WorkerSize

	for i := uint8(0); i < size; i++ {
		go worker(dns, tasks)
	}
}

var (
	once  sync.Once
	tasks chan Task
)

func GetTasks(setting config.Setting) chan Task {
	once.Do(func() {
		dnsSetting := setting.DNS
		bufferSize := dnsSetting.BufferSize
		tasks = make(chan Task, bufferSize)
	})
	return tasks
}

func worker(dns DNS, tasks chan Task) {
	for task := range tasks {
		handler, exists := handlers[task.Request.Method]
		if !exists {
			// TODO: 에러 로깅 추가 [런타임 환경임]
			task.Response <- Response{
				Error:   errors.New(fmt.Sprintf("Unsupported method: %s", task.Request.Method)),
				Records: []Record{},
			}
			continue
		}
		response := handler(task, dns)
		task.Response <- response
	}
}
