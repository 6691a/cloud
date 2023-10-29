package routeros

import (
	"crypto/tls"
	"net/http"
)

type Client struct {
	session *Session
}

func NewClient(apiUrl string, hclient *http.Client, tls *tls.Config, proxyString string) (client *Client, err error) {
	sess, err := NewSession(apiUrl, hclient, tls, proxyString)
	if err != nil {
		return nil, err
	}

	return &Client{
		session: sess,
	}, nil
}

func (c *Client) Login(username, password string) error {
	return c.session.Login(username, password)
}

func (c *Client) GetWhiteList(name, address string) (*WhiteList, error) {
	return c.session.GetWhiteList(name, address)
}

func (c *Client) CreateWhiteList(name, address string) (*WhiteList, error) {
	return c.session.CreateWhiteList(name, address)
}

func (c *Client) CreateDomainList(name, domain string) (*DomainList, error) {
	return c.session.CreateDomainList(name, domain)
}
