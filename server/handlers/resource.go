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

func CSSResource() Handler {
	return &resourceHandler{"/css/{path:.+}", "/css"}
}

func FontResource() Handler {
	return &resourceHandler{"/fonts/{path:.+}", "/fonts"}
}

func FaviconResource() Handler {
	return &resourceHandler{"/{path:[^/]+\\.(?:ico|png|jpg|webmanifest)}", "/images/favicons"}
}

type resourceHandler struct {
	pathTemplate string
	dir string
}

func (h *resourceHandler) Register(r *mux.Router) {
	r.HandleFunc(h.pathTemplate, h.handle)
}

func (h *resourceHandler) handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relativePath := vars["path"]
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
