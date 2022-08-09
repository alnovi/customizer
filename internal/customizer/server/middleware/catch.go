package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"alnovi/customizer/pkg/logger"
)

type CatchMiddleware struct {
}

func NewCatchMiddleware() *CatchMiddleware {
	return &CatchMiddleware{}
}

func (m *CatchMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithError(errors.New("A panic has caught")).
					WithField("recover", r).
					WithField("trace", getBacktrace()).
					Error("internal.server.middleware.catch.CatchMiddleware.Middleware")

				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func getBacktrace() []string {
	return strings.Split(string(debug.Stack()), "\n\t")
}
