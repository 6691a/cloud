package hypervisor

import (
	"github.com/6691a/iac/config"
)

type Hypervisor interface {
	CloneVM(record *CloneRecord) error
	CloneCT(record *CloneRecord) error
	Delete(id int) error
	CreateNetwork(recode *NetworkRecode) error
	SetVMConfig(recode *VMConfigRecode) (interface{}, error)
	SetCTConfig(recode *CTConfigRecode) (interface{}, error)
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
	Id          int
	Name        string
	Description string
}

func NewRecord(id int, name string, description string) *Record {
	return &Record{
		Id:          id,
		Name:        name,
		Description: description,
	}
}

type CloneRecord struct {
	*Record
	NewId  int
	Method Method
	Full   bool
}

func NewCloneRecord(id int, name string, description string, newId int) *CloneRecord {
	return &CloneRecord{
		Record: &Record{
			Id:          id,
			Name:        name,
			Description: description,
		},
		NewId: newId,
		Full:  false,
	}
}

func (r *CloneRecord) SetFull(full bool) {
	r.Full = full
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

type ConfigRecode struct {
	*Record
	core     uint8
	memory   uint16
	netDrive map[uint8]string
	ipCfg    map[uint8]string
	onBoot   bool
	Args     []interface{}
}

type VMConfigRecode struct {
	*ConfigRecode
}

type CTConfigRecode struct {
	*ConfigRecode
}

func newConfigRecode(id int, core uint8, memory uint16, onBoot bool, netDrive map[uint8]string, ipCfg map[uint8]string, args ...interface{}) *ConfigRecode {
	return &ConfigRecode{
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

func NewVmConfigRecode(id int, core uint8, memory uint16, onBoot bool, netDrive map[uint8]string, ipCfg map[uint8]string, args ...interface{}) *VMConfigRecode {
	return &VMConfigRecode{
		ConfigRecode: newConfigRecode(id, core, memory, onBoot, netDrive, ipCfg, args...),
	}
}

func NewCTConfigRecode(id int, core uint8, memory uint16, onBoot bool, netDrive map[uint8]string, ipCfg map[uint8]string, args ...interface{}) *CTConfigRecode {
	return &CTConfigRecode{
		ConfigRecode: newConfigRecode(id, core, memory, onBoot, netDrive, ipCfg, args...),
	}
}
