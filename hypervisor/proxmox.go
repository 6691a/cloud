package hypervisor

import (
	"crypto/tls"
	"fmt"
	"github.com/Telmate/proxmox-api-go/proxmox"
	"strconv"
	"strings"
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

func (p *Proxmox) newCTRef(recode *Record) *proxmox.VmRef {
	vmRef := proxmox.NewVmRef(int(recode.Id))
	vmRef.SetNode(p.Node)
	vmRef.SetVmType("lxc")
	return vmRef
}

func (p *Proxmox) newCloneVMParams(recode *CloneRecord) map[string]interface{} {
	return map[string]interface{}{
		"name":        recode.Name,
		"description": recode.Description,
		"newid":       recode.NewId,
	}
}

func (p *Proxmox) newCloneCTParams(recode *CloneRecord) map[string]interface{} {
	return map[string]interface{}{
		"newid":       recode.NewId,
		"vmid":        strconv.Itoa(recode.Record.Id),
		"hostname":    recode.Name,
		"description": recode.Description,
		"target":      p.Node,
		"node":        p.Node,
		"full":        recode.Full,
	}
}

// TODO: if need to use this function, need to implement it.
func (p *Proxmox) Create() (string, error) {
	//return p.Client.CreateQemuVm()
	return "", nil
}

func (p *Proxmox) CloneVM(recode *CloneRecord) error {
	_, err := p.Client.CloneQemuVm(p.newVmRef(recode.Record), p.newCloneVMParams(recode))

	return err
}

func (p *Proxmox) CloneCT(recode *CloneRecord) error {
	_, err := p.Client.CloneLxcContainer(p.newVmRef(recode.Record), p.newCloneCTParams(recode))
	return err
}

func (p *Proxmox) List() (map[string]interface{}, error) {
	return p.Client.GetVmList()
}

func (p *Proxmox) Delete(id int) error {
	err := p.Shutdown(id)

	if err != nil {
		expectedErrorMsg := fmt.Sprintf("500 CT %d not running", id)
		if !strings.Contains(err.Error(), expectedErrorMsg) {
			return err
		}
	}
	_, err = p.Client.DeleteVm(proxmox.NewVmRef(int(id)))
	return err
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

func (p *Proxmox) Stop(id uint16) error {
	_, err := p.Client.StopVm(proxmox.NewVmRef(int(id)))
	return err
}

func configParams(recode *ConfigRecode) map[string]interface{} {
	params := map[string]interface{}{
		"cores":  recode.core,
		"memory": recode.memory,
		"onboot": recode.onBoot,
	}

	return params
}

func newVMConfigParams(recode *VMConfigRecode) map[string]interface{} {
	params := configParams(recode.ConfigRecode)

	for k, v := range recode.netDrive {
		params[fmt.Sprintf("net%d", k)] = v
	}

	for k, v := range recode.ipCfg {
		params[fmt.Sprintf("ipconfig%d", k)] = v
	}
	return params
}

func (p *Proxmox) SetVMConfig(recode *VMConfigRecode) (interface{}, error) {
	return p.Client.SetVmConfig(p.newVmRef(recode.Record), newVMConfigParams(recode))
}

func newCTConfigParams(recode *CTConfigRecode) map[string]interface{} {
	params := configParams(recode.ConfigRecode)
	//net0:name=eth0,virtio,bridge=vmbr0,ip=192.168.88.250/24,gw=192.168.88.1
	//net0:name=eth0,bridge=vmbr0,ip=192.168.1.100/24,gw=192.168.1.1
	ctn := 0
	for k, deviceConfig := range recode.netDrive {
		ipConfig, ok := recode.ipCfg[k]
		if !ok {
			continue
		}

		name := fmt.Sprintf("name=eth%d,", k)
		params[fmt.Sprintf("net%d", ctn)] = name + deviceConfig + "," + ipConfig
		ctn++
	}
	fmt.Println(params)
	return params
}

func (p *Proxmox) SetCTConfig(recode *CTConfigRecode) (interface{}, error) {
	return p.Client.SetLxcConfig(p.newCTRef(recode.Record), newCTConfigParams(recode))
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
