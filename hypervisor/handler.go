package hypervisor

import (
	"errors"
)

// HandlerFunc Designed to return an interface{} for the sake of flexibility in modifications
type HandlerFunc func(hv Hypervisor, task Task) (interface{}, error)

var handlers = map[string]HandlerFunc{
	"Clone":         handlerClone,
	"Delete":        handlerDelete,
	"CreateNetwork": handlerCreateNetwork,
	"SetVmConfig":   handlerSetVmConfig,
}

func handlerClone(hv Hypervisor, task Task) (interface{}, error) {
	record, ok := task.Request.Record.(*CloneRecord)
	if !ok {
		return nil, errors.New("invalid request type for Clone")
	}
	return nil, hv.Clone(record)
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
	record, ok := task.Request.Record.(*VmConfigRecode)
	if !ok {
		return nil, errors.New("invalid request type for SetVmConfig")
	}
	return hv.SetVmConfig(record)
}
