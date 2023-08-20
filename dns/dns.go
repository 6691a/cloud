package dns

import (
	"errors"
	"fmt"
	"github.com/6691a/iac/config"
	"sync"
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

type Task struct {
	Request  Request
	Response chan Response
}

type Method string

const (
	Get    Method = "GET"
	List   Method = "LIST"
	Create Method = "CREATE"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
)

func CreateWorker(setting config.Setting, tasks chan Task) {
	dns, err := NewDNS(setting)
	if err != nil {
		panic(err)
	}

	dnsSetting := setting.DNS
	size := dnsSetting.WorkerSize

	for i := uint8(0); i < size; i++ {
		go worker(dns, tasks)
	}
}

var (
	once  sync.Once
	tasks chan Task
)

func GetTasks(setting config.Setting) chan Task {
	once.Do(func() {
		dnsSetting := setting.DNS
		bufferSize := dnsSetting.BufferSize
		tasks = make(chan Task, bufferSize)
	})
	return tasks
}

type Request struct {
	Method Method
	Record Record
}

type Response struct {
	Error   error
	Records []Record
}

func worker(dns DNS, tasks chan Task) {
	for task := range tasks {
		var response Response
		switch task.Request.Method {
		case Method(Get):
			record, err := dns.Get(task.Request.Record.SubDomain, task.Request.Record.Type)
			response = Response{
				Error:   err,
				Records: []Record{record},
			}
		case Method(List):
			records, err := dns.List()
			response = Response{
				Error:   err,
				Records: records,
			}
		case Method(Create):
			record, err := dns.Create(task.Request.Record)
			response = Response{
				Error:   err,
				Records: []Record{record},
			}
		case Method(Patch):
			record, err := dns.Patch(task.Request.Record.SubDomain, task.Request.Record.Type, task.Request.Record)
			response = Response{
				Error:   err,
				Records: []Record{record},
			}
		case Method(Delete):
			err := dns.Delete(task.Request.Record.SubDomain, task.Request.Record.Type)
			response = Response{
				Error:   err,
				Records: []Record{},
			}
		}
		task.Response <- response
	}
}
