package routeros

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

type Session struct {
	BaseUrl    string
	HttpClient *http.Client
	user       string
	password   string
}

func NewSession(apiUrl string, hclient *http.Client, tls *tls.Config, proxyString string) (*Session, error) {
	if hclient == nil {
		if proxyString != "" {
			proxyURL, err := url.ParseRequestURI(proxyString)
			if err != nil {
				return nil, err
			}
			if _, _, err := net.SplitHostPort(proxyURL.Host); err != nil {
				return nil, err
			}
			tp := &http.Transport{
				TLSClientConfig:    tls,
				DisableCompression: true,
				Proxy:              http.ProxyURL(proxyURL),
			}
			hclient = &http.Client{Transport: tp}
		} else {
			tp := &http.Transport{
				TLSClientConfig:    tls,
				DisableCompression: true,
				Proxy:              nil,
			}
			hclient = &http.Client{Transport: tp}
		}
	}

	return &Session{
		HttpClient: hclient,
		BaseUrl:    apiUrl,
	}, nil
}

func (s *Session) NewRequest(method, url string, params *url.Values, headers *http.Header, body io.Reader) (*http.Request, error) {
	url = s.createParams(fmt.Sprintf("%s/rest/%s", s.BaseUrl, url), nil)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	// TODO: Basic Auth
	req.SetBasicAuth(s.user, s.password)
	return req, nil
}

func (s *Session) Login(username, password string) error {
	s.user = username
	s.password = password

	req, err := s.NewRequest("GET", "system/resource", nil, nil, nil)
	if err != nil {
		return err
	}
	res, err := s.HttpClient.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("Login failed: %s", bodyString)
	}
	return nil
}

func (s *Session) createParams(url string, params *url.Values) string {
	if params != nil {
		return fmt.Sprintf("%s?%s", url, params.Encode())
	}
	return url
}
