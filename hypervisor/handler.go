package hypervisor

import "errors"

type HypervisorHandlerFunc func(hv Hypervisor, task Task) error

var handlers = map[string]HypervisorHandlerFunc{
	"Clone":         handlerClone,
	"Delete":        handlerDelete,
	"CreateNetwork": handlerCreateNetwork,
	"SetVmConfig":   handlerSetVmConfig,
}

func handlerClone(hv Hypervisor, task Task) error {
	record, ok := task.Request.Record.(*CloneRecord)
	if !ok {
		return errors.New("invalid request type for Clone")
	}
	return hv.Clone(record)
}

func handlerDelete(hv Hypervisor, task Task) error {
	id, ok := task.Request.Record.(uint16)
	if !ok {
		return errors.New("invalid request type for Delete")
	}
	return hv.Delete(id)
}

func handlerCreateNetwork(hv Hypervisor, task Task) error {
	record, ok := task.Request.Record.(*NetworkRecode)
	if !ok {
		return errors.New("invalid request type for CreateNetwork")
	}
	return hv.CreateNetwork(record)
}

func handlerSetVmConfig(hv Hypervisor, task Task) error {
	record, ok := task.Request.Record.(*VmConfigRecode)
	if !ok {
		return errors.New("invalid request type for SetVmConfig")
	}
	return hv.SetVmConfig(record)
}
