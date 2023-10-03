package routeros

import "crypto/tls"

type RouterOS struct {
	Client *Client
}

func NewRouterOS(url string, username string, password string) (*RouterOS, error) {
	tls := &tls.Config{InsecureSkipVerify: true}
	client, err := NewClient(url, nil, tls, "")
	if err != nil {
		return nil, err
	}
	return &RouterOS{
		Client: client,
	}, nil
}

func (r *RouterOS) Login(username string, password string) error {
	return r.Client.Login(username, password)
}
