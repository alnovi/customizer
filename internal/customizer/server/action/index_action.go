package action

import "net/http"

type IndexHandler struct {
	root string
}

func NewIndexHandler(root string) *IndexHandler {
	return &IndexHandler{
		root: root,
	}
}

func (h *IndexHandler) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	write.WriteHeader(http.StatusOK)
}
