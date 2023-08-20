package hypervisor

import "github.com/6691a/iac/config"

type Hypervisor interface {
	Clone(record *CloneRecord) (string, error)
	Delete(id int) (string, error)
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
	id          int
	name        string
	description string
}

func NewRecord(id int, name string, description string) *Record {
	return &Record{
		id:          id,
		name:        name,
		description: description,
	}
}

type CloneRecord struct {
	Record
	newId int
}

func NewCloneRecord(id int, name string, description string, newId int) *CloneRecord {
	return &CloneRecord{
		Record: Record{
			id:          id,
			name:        name,
			description: description,
		},
		newId: newId,
	}
}
