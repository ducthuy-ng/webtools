package webtools

import (
	"fmt"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ViteIntegrationConfigs struct {
	isDevEnvironment bool
	ViteProxyAddress string
	ViteBuildPath    string
	PathPrefix       string
}

func (configs ViteIntegrationConfigs) SetIsDevEnvironment(yes bool) ViteIntegrationConfigs {
	configs.isDevEnvironment = yes
	return configs
}

func NewViteIntegrationConfigs(vitebuildPath string) ViteIntegrationConfigs {
	return ViteIntegrationConfigs{
		isDevEnvironment: false,
		ViteProxyAddress: "http://localhost:5173/",
		ViteBuildPath:    vitebuildPath,
		PathPrefix:       "/",
	}
}

func ApplyViteIntegration(e *echo.Echo, configs ViteIntegrationConfigs) error {
	if configs.isDevEnvironment {
		e.Logger.Info(configs.ViteProxyAddress)
		err := setupDevProxy(e, configs)
		if err != nil {
			return err
		}
	}

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   configs.ViteBuildPath,
		HTML5:  true,
		Index:  "index.html",
		Browse: true}))
	return nil
}

func setupDevProxy(e *echo.Echo, configs ViteIntegrationConfigs) error {
	url, err := url.Parse(configs.ViteProxyAddress)
	if err != nil {
		return fmt.Errorf("failed to parse Vite proxy address: %v", err)
	}
	proxyTargets := []*middleware.ProxyTarget{{Name: "Vite Dev", URL: url}}
	balancer := middleware.NewRoundRobinBalancer(proxyTargets)
	e.Use(middleware.Proxy(balancer))
	return nil
}
