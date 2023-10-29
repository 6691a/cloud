package hypervisor

import (
	"github.com/6691a/iac/config"
)

type Hypervisor interface {
	Clone(record *CloneRecord) error
	Delete(id uint16) error
	CreateNetwork(recode *NetworkRecode) error
	SetVmConfig(recode *VmConfigRecode) (interface{}, error)
}

func NewHypervisor(setting config.Setting) (Hypervisor, error) {
	hvSetting := setting.Hypervisor
	switch hvSetting.Service {
	case "proxmox":
		proxmoxSetting := hvSetting.Proxmox
		return NewProxmox(proxmoxSetting.Url, proxmoxSetting.Node, proxmoxSetting.User, proxmoxSetting.Token, 300)
	}
	return nil, nil
}

type Record struct {
	Id          uint16
	Name        string
	Description string
}

func NewRecord(id uint16, name string, description string) *Record {
	return &Record{
		Id:          id,
		Name:        name,
		Description: description,
	}
}

type CloneRecord struct {
	*Record
	NewId  uint16
	Method Method
}

func NewCloneRecord(id uint16, name string, description string, newId uint16) *CloneRecord {
	return &CloneRecord{
		Record: &Record{
			Id:          id,
			Name:        name,
			Description: description,
		},
		NewId:  newId,
		Method: Clone,
	}
}

type NetworkType string

const (
	Bridge     NetworkType = "bridge"
	Bond       NetworkType = "bond"
	VLAN       NetworkType = "vlan"
	OVSBridge  NetworkType = "OVSBridge"
	OVSBond    NetworkType = "OVSBond"
	OVSIntPort NetworkType = "OVSIntPort"
)

type NetworkRecode struct {
	*Record
	Type_     NetworkType
	AutoStart bool
	Cidr      string
	Args      []interface{}
}

func NewNetworkRecode(name, cidr, description string, type_ NetworkType, autoStart bool, args ...interface{}) *NetworkRecode {
	return &NetworkRecode{
		Record: &Record{
			Id:          0,
			Name:        name,
			Description: description,
		},
		Type_:     type_,
		AutoStart: autoStart,
		Cidr:      cidr,
		Args:      args,
	}
}

type VmConfigRecode struct {
	*Record
	core     uint8
	memory   uint16
	netDrive map[uint8]string
	ipCfg    map[uint8]string
	onBoot   bool
	Args     []interface{}
}

func NewVmConfigRecode(id uint16, core uint8, memory uint16, onBoot bool, netDrive map[uint8]string, ipCfg map[uint8]string, args ...interface{}) *VmConfigRecode {
	return &VmConfigRecode{
		Record: &Record{
			Id:          id,
			Name:        "",
			Description: "",
		},
		core:     core,
		memory:   memory,
		netDrive: netDrive,
		ipCfg:    ipCfg,
		onBoot:   onBoot,
		Args:     args,
	}
}
