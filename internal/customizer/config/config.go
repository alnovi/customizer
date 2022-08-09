package config

import (
	"alnovi/customizer/internal/customizer/server"
	"alnovi/customizer/pkg/storage/drivers"
	"time"
)

type Config struct {
	HttpServer   server.HttpServerConfig
	MongoConfig  drivers.MongoConfig
	SentryConfig SentryConfig
}

type SentryConfig struct {
	Url     string
	Debug   bool
	Timeout time.Duration
}
