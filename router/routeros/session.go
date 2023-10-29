package routeros

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/6691a/iac/config"
	"go.uber.org/zap"
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

func (s *Session) createParams(url string, params *url.Values) string {
	if params != nil {
		return fmt.Sprintf("%s?%s", url, params.Encode())
	}
	return url
}

func toJson(body map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func validateResponse(res *http.Response) error {
	logger := config.GetLogger("default")
	logger.Error("request failed", zap.String("body", res.Status))
	if res.StatusCode/100 != 2 {
		logger := config.GetLogger("default")
		bodyBytes, _ := io.ReadAll(res.Body)
		bodyString := string(bodyBytes)
		logger.Error("request failed", zap.String("body", bodyString))
		return fmt.Errorf("request failed: %s", bodyString)
	}
	return nil
}

func (s *Session) JsonRequest(method, url string, params *url.Values, headers *http.Header, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := s.NewRequest(
		method, url, params, headers, &jsonData,
	)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *Session) Do(req *http.Request, result interface{}) error {
	res, err := s.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = validateResponse(res)
	if err != nil {
		return err
	}

	if result != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) NewRequest(method, url string, params *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	url = s.createParams(fmt.Sprintf("%s/rest/%s", s.BaseUrl, url), nil)

	var buf io.Reader
	if body != nil {
		buf = bytes.NewReader(*body)
	}

	req, err := http.NewRequest(method, url, buf)
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

	req, err := s.JsonRequest(
		"GET",
		"system/resource",
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	err = s.Do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) GetWhiteList(name, address string) (*WhiteList, error) {
	req, err := s.JsonRequest(
		"GET",
		"ip/firewall/address-list",
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var entries []WhiteList
	err = s.Do(req, &entries)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.List == name && entry.Address == address {
			return &entry, nil
		}
	}

	return nil, nil
}

func (s *Session) CreateWhiteList(name, address string) (*WhiteList, error) {
	req, err := s.JsonRequest(
		"PUT",
		"ip/firewall/address-list",
		nil,
		nil,
		map[string]interface{}{
			"list":    name,
			"address": address,
		},
	)
	if err != nil {
		return nil, err
	}

	var entries WhiteList
	err = s.Do(req, &entries)
	if err != nil {
		return nil, err
	}

	return &entries, nil
}

func (s *Session) DeleteWhiteList(id string) error {
	req, err := s.JsonRequest(
		"DELETE",
		fmt.Sprintf("ip/firewall/address-list/%s", id),
		nil,
		nil,
		nil,
	)

	if err != nil {
		return err
	}

	err = s.Do(req, nil)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) CreateDomainList(name, domain string) (*DomainList, error) {
	req, err := s.JsonRequest(
		"PUT",
		"ip/firewall/layer7-protocol",
		nil,
		nil,
		map[string]interface{}{
			"name":   name,
			"regexp": domain,
		},
	)

	if err != nil {
		return nil, err
	}

	var entries DomainList

	err = s.Do(req, &entries)
	if err != nil {
		return nil, err
	}
	return &entries, nil
}
