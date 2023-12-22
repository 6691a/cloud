package hypervisor

import (
	"errors"
)

// HandlerFunc Designed to return an interface{} for the sake of flexibility in modifications
type HandlerFunc func(hv Hypervisor, task Task) (interface{}, error)

var handlers = map[string]HandlerFunc{
	"CloneCT":       handlerCloneCT,
	"CloneVM":       handlerCloneVM,
	"Delete":        handlerDelete,
	"CreateNetwork": handlerCreateNetwork,
	"SetVMConfig":   handlerSetVmConfig,
	"SetCTConfig":   handlerSetCTConfig,
}

func handlerCloneCT(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*CloneRecord)

	if !ok {
		return nil, errors.New("invalid request type for Clone")
	}

	return nil, hv.CloneCT(record)
}

func handlerCloneVM(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*CloneRecord)

	if !ok {
		return nil, errors.New("invalid request type for Clone")
	}

	return nil, hv.CloneVM(record)
}

func handlerDelete(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*Record)
	if !ok {
		return nil, errors.New("invalid request type for Delete")
	}
	id := record.Id

	return nil, hv.Delete(id)
}

func handlerCreateNetwork(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*NetworkRecode)
	if !ok {
		return nil, errors.New("invalid request type for CreateNetwork")
	}
	return nil, hv.CreateNetwork(record)
}

func handlerSetVmConfig(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*VMConfigRecode)
	if !ok {
		return nil, errors.New("invalid request type for SetVmConfig")
	}
	return hv.SetVMConfig(record)
}

func handlerSetCTConfig(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*CTConfigRecode)
	if !ok {
		return nil, errors.New("invalid request type for SetCTConfig")
	}
	return hv.SetCTConfig(record)
}
