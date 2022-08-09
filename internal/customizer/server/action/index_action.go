package action

import (
	"net/http"
	"os"
	"path/filepath"
)

type IndexHandler struct {
	root string
}

func NewIndexHandler(root string) *IndexHandler {
	return &IndexHandler{
		root: root,
	}
}

func (h *IndexHandler) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(request.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.root, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(write, request, filepath.Join(h.root, "index.html"))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(write, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.root)).ServeHTTP(write, request)
}
