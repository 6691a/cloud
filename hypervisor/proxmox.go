package hypervisor

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
)

type Proxmox struct {
	Client *proxmox.Client
	Node   string
}

func NewProxmox(url, node, user, token string, timeOut uint16) (*Proxmox, error) {
	tls := &tls.Config{InsecureSkipVerify: true}
	client, err := proxmox.NewClient(url, nil, "", tls, "", int(timeOut))

	if err != nil {
		return nil, err
	}

	client.SetAPIToken(user, token)

	return &Proxmox{
		Client: client,
		Node:   node,
	}, nil
}

func (p *Proxmox) newVmRef(recode *Record) *proxmox.VmRef {
	vmRef := proxmox.NewVmRef(int(recode.Id))
	vmRef.SetNode(p.Node)
	vmRef.SetVmType("qemu")
	return vmRef
}

func (p *Proxmox) newCloneParams(recode *CloneRecord) map[string]interface{} {
	return map[string]interface{}{
		"name":        recode.Name,
		"description": recode.Description,
		"newid":       recode.NewId,
	}
}

// TODO: if need to use this function, need to implement it.
func (p *Proxmox) Create() (string, error) {
	//return p.Client.CreateQemuVm()
	return "", nil
}

func (p *Proxmox) Clone(recode *CloneRecord) error {
	_, err := p.Client.CloneQemuVm(p.newVmRef(recode.Record), p.newCloneParams(recode))
	return err
}

func (p *Proxmox) List() (map[string]interface{}, error) {
	return p.Client.GetVmList()
}

func (p *Proxmox) Delete(id uint16) error {
	err := p.Shutdown(id)

	if err != nil {
		return err
	}

	_, err = p.Client.DeleteVm(proxmox.NewVmRef(int(id)))
	return err
}

func (p *Proxmox) Shutdown(id uint16) error {
	_, err := p.Client.ShutdownVm(proxmox.NewVmRef(int(id)))
	return err
}

func (p *Proxmox) Reboot(id uint16) error {
	_, err := p.Client.ResetVm(proxmox.NewVmRef(int(id)))
	return err
}

func (p *Proxmox) Start(id uint16) error {
	_, err := p.Client.StartVm(proxmox.NewVmRef(int(id)))
	return err
}

func (p *Proxmox) Stop(id uint16) error {
	_, err := p.Client.StopVm(proxmox.NewVmRef(int(id)))
	return err
}

func newVmConfigParams(recode *VmConfigRecode) map[string]interface{} {
	params := map[string]interface{}{
		"cores":  recode.core,
		"memory": recode.memory,
		"onboot": recode.onBoot,
	}

	for k, v := range recode.netDrive {
		params[fmt.Sprintf("net%d", k)] = v
	}

	for k, v := range recode.ipCfg {
		params[fmt.Sprintf("ipconfig%d", k)] = v
	}

	return params
}

func (p *Proxmox) SetVmConfig(recode *VmConfigRecode) error {
	_, err := p.Client.SetVmConfig(p.newVmRef(recode.Record), newVmConfigParams(recode))
	return err
}

// ======================== Node Network ========================
func newNetworkParams(recode *NetworkRecode) map[string]interface{} {
	params := map[string]interface{}{
		"iface":     recode.Name,
		"type":      recode.Type_,
		"comments":  recode.Description,
		"cidr":      recode.Cidr,
		"autostart": recode.AutoStart,
	}
	for _, arg := range recode.Args {
		if m, ok := arg.(map[string]interface{}); ok {
			for k, v := range m {
				params[k] = v
			}
		}
	}
	return params
}

func (p *Proxmox) CreateNetwork(recode *NetworkRecode) error {
	_, err := p.Client.CreateNetwork(p.Node, newNetworkParams(recode))

	if err != nil {
		return err
	}

	return p.ApplyNetwork()
}

func (p *Proxmox) ApplyNetwork() error {
	_, err := p.Client.ApplyNetwork(p.Node)
	return err
}

// TODO: if need to use this function, need to implement it.
func (p *Proxmox) DeleteNetwork(recode Record) error {
	//_, err := p.Client.DeleteNetwork()
	return nil
}
