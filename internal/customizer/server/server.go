package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/braintree/manners"
)

type HttpServerConfig struct {
	Host    string
	Port    uint
	Timeout uint
}

type HttpServerAdapter struct {
	server *manners.GracefulServer
}

func NewHttpServerAdapter(config HttpServerConfig) *HttpServerAdapter {
	timeout := time.Duration(config.Timeout) * time.Second

	httpServer := http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	server := manners.NewWithServer(&httpServer)

	return &HttpServerAdapter{
		server: server,
	}
}

func (hsa *HttpServerAdapter) SetHandler(handler http.Handler) {
	hsa.server.Handler = handler
}

func (hsa *HttpServerAdapter) ListenAndServe() error {
	return hsa.server.ListenAndServe()
}

func (hsa *HttpServerAdapter) Close() bool {
	return hsa.server.Close()
}
