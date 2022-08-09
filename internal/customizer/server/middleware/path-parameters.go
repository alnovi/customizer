package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type PathParametersMiddleware struct {
}

func NewPathParametersMiddleware() *PathParametersMiddleware {
	return &PathParametersMiddleware{}
}

func (p PathParametersMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parameters := mux.Vars(r)

		ctx := context.WithValue(r.Context(), "parameters", parameters)

		next.ServeHTTP(w, r.WithContext(ctx))

		return
	})
}
