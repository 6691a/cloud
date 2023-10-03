package router

import (
	"errors"
	"fmt"
	"github.com/6691a/iac/config"
	"github.com/6691a/iac/router/routeros"
)

type Router interface {
	Login(username string, password string) error
}

func NewRouter(setting config.Setting) (Router, error) {
	rtSetting := setting.Router
	switch rtSetting.Service {
	case "routeros":
		routerOSSetting := rtSetting.RouterOS
		return routeros.NewRouterOS(routerOSSetting.Url, routerOSSetting.User, routerOSSetting.Password)
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported Router service: %s", rtSetting.Service))
	}
}
