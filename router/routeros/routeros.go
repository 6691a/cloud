package routeros

import (
	"crypto/tls"
	"github.com/6691a/iac/config"
	"go.uber.org/zap"
)

type RouterOS struct {
	Client *Client
	logger *zap.Logger
}

// address-list structure
type WhiteList struct {
	ID           string `json:".id"`
	Address      string `json:"address"`
	CreationTime string `json:"creation-time"`
	Disabled     string `json:"disabled"`
	Dynamic      string `json:"dynamic"`
	List         string `json:"list"`
}

// layer7-protocol structure
type DomainList struct {
	ID     string `json:".id"`
	Name   string `json:"name"`
	Domain string `json:"regexp"`
}

func NewRouterOS(url string, username string, password string) (*RouterOS, error) {
	tls := &tls.Config{InsecureSkipVerify: true}
	client, err := NewClient(url, nil, tls, "")

	if err != nil {
		return nil, err
	}
	router := &RouterOS{
		Client: client,
		logger: config.GetLogger("default"),
	}

	// TODO: basic auth login check change to different approach (if supported in RouterOS)
	if err := router.Login(username, password); err != nil {
		router.logger.Panic("RouterOS login failed", zap.Error(err))
		return nil, err
	}
	return router, nil
}

func (r *RouterOS) Login(username string, password string) error {
	return r.Client.Login(username, password)
}

func (r *RouterOS) GetWhiteList(name, address string) (*WhiteList, error) {
	return r.Client.GetWhiteList(name, address)
}
func (r *RouterOS) CreateWhiteList(name, address string) (*WhiteList, error) {
	return r.Client.CreateWhiteList(name, address)
}

func (r *RouterOS) CreateDomainList(name, domain string) (*DomainList, error) {
	return r.Client.CreateDomainList(name, domain)
}

func (r *RouterOS) DeleteWhiteList() {}

func (r *RouterOS) AddRoute() {
}

func (r *RouterOS) DeleteRoute() {}

func (r *RouterOS) AddMangle() {}

func (r *RouterOS) DeleteMangle() {}

func (r *RouterOS) UpdateWhiteList() {}

func (r *RouterOS) AddFilter() {}

func (r *RouterOS) DeleteFilter() {}
