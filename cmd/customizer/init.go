package main

import (
	internalConfig "alnovi/customizer/internal/customizer/config"
	"alnovi/customizer/internal/customizer/server"
	configWrapper "alnovi/customizer/pkg/config"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/storage/drivers"
	"time"

	"github.com/sirupsen/logrus"
)

func initConfig() internalConfig.Config {
	config := configWrapper.All
	storageConfig := configWrapper.Storage

	return internalConfig.Config{
		HttpServer: server.HttpServerConfig{
			Host:    config.HttpServer.Host,
			Port:    uint(config.HttpServer.Port),
			Timeout: uint(config.HttpServer.Timeout),
		},
		MongoConfig: drivers.MongoConfig{
			Host:     storageConfig.Mongo.Host,
			Port:     uint(storageConfig.Mongo.Port),
			Database: storageConfig.Mongo.Database,
			User:     storageConfig.Mongo.User,
			Password: storageConfig.Mongo.Password,
			Timeout:  storageConfig.Mongo.Timeout,
		},
		SentryConfig: internalConfig.SentryConfig{
			Url:     config.Sentry.Url,
			Debug:   config.Sentry.Debug,
			Timeout: time.Duration(config.Sentry.Timeout) * time.Second,
		},
	}
}

func initLogger() {
	logger.SetFormatter(transformLogFormat(configWrapper.Log.Format))
	logger.SetLevel(transformLogLevel(configWrapper.Log.Level))
}

func transformLogLevel(level string) logrus.Level {
	switch level {
	case "critical":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warning":
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.WarnLevel
	}

	return logrus.WarnLevel
}

func transformLogFormat(format string) logrus.Formatter {
	switch format {
	case "text":
		return &logrus.TextFormatter{}
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}
