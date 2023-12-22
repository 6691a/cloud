package hypervisor

import (
	"errors"
	"github.com/6691a/iac/config"
	"sync"
)

type Method string

const (
	CloneCT       Method = "CloneCT"
	CloneVM       Method = "CloneVM"
	Delete        Method = "Delete"
	CreateNetwork Method = "CreateNetwork"
	SetVmConfig   Method = "SetVMConfig"
	SetCTConfig   Method = "SetCTConfig"
)

type Request struct {
	Method Method
	Record interface{}
}

type Response struct {
	Error   error
	Records interface{}
}

type Task struct {
	Request  Request
	Response chan Response
}

var (
	once  sync.Once
	tasks chan Task
)

func GetTasks(setting config.Setting) chan Task {
	once.Do(func() {
		dnsSetting := setting.Hypervisor
		bufferSize := dnsSetting.BufferSize
		tasks = make(chan Task, bufferSize)
	})
	return tasks
}

func worker(hv Hypervisor, tasks chan Task) {
	for task := range tasks {
		method := string(task.Request.Method)
		handler, exists := handlers[method]
		if !exists {
			// TODO: log 작성
			task.Response <- Response{
				Error: errors.New("unsupported method: " + method),
			}
			continue
		}
		records, err := handler(hv, task)
		task.Response <- Response{
			Error:   err,
			Records: records,
		}
	}
}

func CreateWorker(setting config.Setting, tasks chan Task) {
	hv, err := NewHypervisor(setting)
	if err != nil {
		panic(err)
	}

	hvSetting := setting.Hypervisor
	size := hvSetting.WorkerSize

	for i := uint8(0); i < size; i++ {
		go worker(hv, tasks)
	}
}
