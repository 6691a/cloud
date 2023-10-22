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
	router.logger.Error("RouterOS login failed", zap.Error(err))
	// TODO: basic auth login check
	if err := router.Login(username, password); err != nil {
		router.logger.Panic("RouterOS login failed", zap.Error(err))
		return nil, err
	}
	return router, nil
}

func (r *RouterOS) Login(username string, password string) error {
	return r.Client.Login(username, password)
}
