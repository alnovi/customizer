package main

import (
	internalConfig "alnovi/customizer/internal/customizer/config"
	"alnovi/customizer/internal/customizer/server"
	configWrapper "alnovi/customizer/pkg/config"
)

func initConfig() internalConfig.Config {
	config := configWrapper.All

	return internalConfig.Config{
		HttpServer: server.HttpServerConfig{
			Host:    config.HttpServer.Host,
			Port:    uint(config.HttpServer.Port),
			Timeout: uint(config.HttpServer.Timeout),
		},
	}
}
