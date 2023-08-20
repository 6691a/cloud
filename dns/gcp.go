package dns

import (
	"context"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

type GCP struct {
	service     *dns.Service
	projectId   string
	managedZone string
	domainName  string
}

func NewGCP(projectID, managedZone, credentialsPath, domainName string) (*GCP, error) {
	ctx := context.Background()
	dnsService, err := dns.NewService(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}
	return &GCP{
		service:     dnsService,
		projectId:   projectID,
		managedZone: managedZone,
		domainName:  domainName,
	}, nil
}

func (g *GCP) Get(subDomain string, type_ Type) (Record, error) {
	rcd, err := g.service.ResourceRecordSets.Get(g.projectId, g.managedZone, subDomain+g.domainName, string(type_)).Do()
	if err != nil {
		return Record{}, err
	}

	return Record{
		SubDomain: rcd.Name,
		Type:      Type(rcd.Type),
		Ttl:       rcd.Ttl,
		RtDatas:   rcd.Rrdatas,
	}, nil
}

func (g *GCP) List() ([]Record, error) {
	rcds, err := g.service.ResourceRecordSets.List(g.projectId, g.managedZone).Do()
	if err != nil {
		return nil, err
	}

	var result []Record

	for _, record := range rcds.Rrsets {
		result = append(result, Record{
			SubDomain: record.Name,
			Type:      Type(record.Type),
			Ttl:       record.Ttl,
			RtDatas:   record.Rrdatas,
		})
	}

	return result, nil
}

func (g *GCP) Create(rcd Record) (Record, error) {
	recodeSet := &dns.ResourceRecordSet{
		Name:    rcd.SubDomain + g.domainName,
		Rrdatas: rcd.RtDatas,
		Ttl:     rcd.Ttl,
		Type:    string(rcd.Type),
	}

	resRcd, err := g.service.ResourceRecordSets.Create(g.projectId, g.managedZone, recodeSet).Do()
	if err != nil {
		return Record{}, err
	}

	return Record{
		SubDomain: resRcd.Name,
		Type:      Type(resRcd.Type),
		Ttl:       resRcd.Ttl,
		RtDatas:   resRcd.Rrdatas,
	}, nil
}

func (g *GCP) Patch(subDomain string, type_ Type, rc Record) (Record, error) {
	recodeSet := &dns.ResourceRecordSet{
		Name:    rc.SubDomain,
		Rrdatas: rc.RtDatas,
		Ttl:     rc.Ttl,
		Type:    string(rc.Type),
	}

	rcd, err := g.service.ResourceRecordSets.Patch(g.projectId, g.managedZone, subDomain+g.domainName, string(type_), recodeSet).Do()
	if err != nil {
		return Record{}, err
	}

	return Record{
		SubDomain: rcd.Name,
		Type:      Type(rcd.Type),
		Ttl:       rcd.Ttl,
		RtDatas:   rcd.Rrdatas,
	}, nil
}

func (g *GCP) Delete(subDomain string, type_ Type) error {
	_, err := g.service.ResourceRecordSets.Delete(g.projectId, g.managedZone, subDomain+g.domainName, string(type_)).Do()
	return err
}
