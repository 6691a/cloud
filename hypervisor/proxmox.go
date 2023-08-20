package hypervisor

import (
	"crypto/tls"
	"github.com/Telmate/proxmox-api-go/proxmox"
)

type Proxmox struct {
	Client *proxmox.Client
	Node   string
}

func NewProxmox(url, node, user, token string, timeOut int) (*Proxmox, error) {
	tls := &tls.Config{InsecureSkipVerify: true}
	client, err := proxmox.NewClient(url, nil, "", tls, "", timeOut)

	if err != nil {
		return nil, err
	}

	client.SetAPIToken(user, token)

	return &Proxmox{
		Client: client,
		Node:   node,
	}, nil
}

func (p *Proxmox) newVmRef(recode *CloneRecord) *proxmox.VmRef {
	vmRef := proxmox.NewVmRef(recode.id)
	vmRef.SetNode(p.Node)
	vmRef.SetVmType("qemu")
	return vmRef
}

func (p *Proxmox) newParams(recode Record) map[string]interface{} {
	return map[string]interface{}{
		"name":        recode.name,
		"description": recode.description,
	}
}

func (p *Proxmox) newCloneParams(recode *CloneRecord) map[string]interface{} {
	params := p.newParams(recode.Record)
	params["newid"] = recode.newId
	return params
}

// TODO: if need to use this function, need to implement it.
func (p *Proxmox) Create() (string, error) {
	//return p.Client.CreateQemuVm()
	return "", nil
}

func (p *Proxmox) Clone(recode *CloneRecord) (string, error) {
	vmRef := p.newVmRef(recode)

	status, err := p.Client.CloneQemuVm(vmRef, p.newCloneParams(recode))
	if err != nil {
		return "", err
	}

	err = p.Start(recode.newId)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (p *Proxmox) List() (map[string]interface{}, error) {
	return p.Client.GetVmList()
}

func (p *Proxmox) Delete(id int) (string, error) {
	err := p.Shutdown(id)

	if err != nil {
		return "", err
	}

	return p.Client.DeleteVm(proxmox.NewVmRef(id))
}

func (p *Proxmox) Shutdown(id int) error {
	_, err := p.Client.ShutdownVm(proxmox.NewVmRef(id))
	return err
}

func (p *Proxmox) Reboot(id int) error {
	_, err := p.Client.ResetVm(proxmox.NewVmRef(id))
	return err
}

func (p *Proxmox) Start(id int) error {
	_, err := p.Client.StartVm(proxmox.NewVmRef(id))
	return err
}

func (p *Proxmox) Stop(id int) error {
	_, err := p.Client.StopVm(proxmox.NewVmRef(id))
	return err
}
