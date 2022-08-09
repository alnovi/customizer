package action

import (
	"alnovi/customizer/internal/customizer/storage"
	"fmt"
	"net/http"
)

type HealthCheckAction struct {
	storage storage.HealthChecker
	cache   storage.HealthChecker
}

func NewHealthCheckAction(storage storage.HealthChecker, cache storage.HealthChecker) HealthCheckAction {
	return HealthCheckAction{
		storage: storage,
		cache:   cache,
	}
}

func (h HealthCheckAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	isStorageAvailable := h.storage.HealthCheck()
	isCacheAvailable := h.cache.HealthCheck()

	_, _ = fmt.Fprintf(writer, "customizer{customizer} 1\n")
	_, _ = fmt.Fprintf(writer, "customizer{storage} %d\n", isStorageAvailable)
	_, _ = fmt.Fprintf(writer, "customizer{cache} %d\n", isCacheAvailable)
}
