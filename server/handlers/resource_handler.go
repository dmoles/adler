package handlers

import (
	"github.com/dmoles/adler/server/resources"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type resourceHandler struct {
	dir     string
	varname string
}

// TODO: optionally log errors?
func (h *resourceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relativePath := vars[h.varname]
	if relativePath == "" || strings.Contains(relativePath, "..") {
		http.NotFound(w, r)
		return
	}
	relativePathClean := path.Clean(relativePath)
	resourcePath := path.Join(h.dir, relativePathClean)

	resource, err := resources.Get(resourcePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentType := resource.ContentType()
	w.Header().Add("Content-Type", contentType)

	size := resource.Size()
	w.Header().Add("Content-Length", strconv.FormatInt(size, 10))

	n, err := resource.Copy(w)
	if err != nil {
		log.Printf("Error serving %#v: %v", resourcePath, err)
	}
	if n != size {
		log.Printf("Wrote wrong number of bytes for %#v: expected %d, was %d", resourcePath, size, n)
	}
}
