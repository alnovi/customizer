package middleware

import "net/http"

type ResponseHeadersMiddleware struct {
}

func NewResponseHeadersMiddleware() *ResponseHeadersMiddleware {
	return &ResponseHeadersMiddleware{}
}

func (m *ResponseHeadersMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)

		return
	})
}
