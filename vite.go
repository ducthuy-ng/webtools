package webtools

import (
	"fmt"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ViteIntegrationConfigs struct {
	isDevEnvironment bool
	viteProxyAddress string
	viteBuildPath    string
	pathPrefix       string
}

func (configs ViteIntegrationConfigs) SetIsDevEnvironment(yes bool) ViteIntegrationConfigs {
	configs.isDevEnvironment = yes
	return configs
}

func NewViteIntegrationConfigs(vitebuildPath string) ViteIntegrationConfigs {
	return ViteIntegrationConfigs{
		isDevEnvironment: false,
		viteProxyAddress: "http://localhost:5173/",
		viteBuildPath:    vitebuildPath,
		pathPrefix:       "/",
	}
}

func ApplyViteIntegration(e *echo.Echo, configs ViteIntegrationConfigs) error {
	if configs.isDevEnvironment {
		e.Logger.Info(configs.viteProxyAddress)
		err := setupDevProxy(e, configs)
		if err != nil {
			return err
		}
	}

	e.Static(configs.pathPrefix, configs.viteBuildPath)
	return nil
}

func setupDevProxy(e *echo.Echo, configs ViteIntegrationConfigs) error {
	url, err := url.Parse(configs.viteProxyAddress)
	if err != nil {
		return fmt.Errorf("failed to parse Vite proxy address: %v", err)
	}
	proxyTargets := []*middleware.ProxyTarget{{Name: "Vite Dev", URL: url}}
	balancer := middleware.NewRoundRobinBalancer(proxyTargets)
	e.Use(middleware.Proxy(balancer))
	return nil
}
