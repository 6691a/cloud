package dns

import (
	"errors"
	"fmt"
	"github.com/6691a/iac/config"
)

type Type string

const (
	A     Type = "A"
	AAAA  Type = "AAAA"
	CNAME Type = "CNAME"
	MX    Type = "MX"
	NS    Type = "NS"
)

type Record struct {
	SubDomain string
	Type      Type
	Ttl       int64
	RtDatas   []string // Routing Datas
}

func NewRecord(subDomain string, type_ Type, ttl int64, rtDatas []string) *Record {
	return &Record{
		SubDomain: subDomain,
		Type:      type_,
		Ttl:       ttl,
		RtDatas:   rtDatas,
	}
}

type DNS interface {
	Get(subDomain string, type_ Type) (Record, error)
	List() ([]Record, error)
	Create(rcd Record) (Record, error)
	Patch(subDomain string, type_ Type, rcd Record) (Record, error)
	Delete(subDomain string, type_ Type) error
}

func NewDNS(setting config.Setting) (DNS, error) {
	dnsSetting := setting.DNS
	domain := "." + dnsSetting.Domain + "."

	switch dnsSetting.Service {
	case "gcp":
		gcpSetting := dnsSetting.GCP
		return NewGCP(gcpSetting.ProjectId, gcpSetting.ManagedZone, gcpSetting.CredentialPath, domain)
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported DNS service: %s", dnsSetting.Service))
	}
}
