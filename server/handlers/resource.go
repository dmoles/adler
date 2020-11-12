package handlers

import (
	"github.com/dmoles/adler/server/resources"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func FontResource() Handler {
	return &resourceHandler{"/fonts/{path:.+}", "/fonts"}
}

func FaviconResource() Handler {
	return &resourceHandler{"/{path:[^/]+\\.(?:ico|png|jpg|webmanifest)}", "/images/favicons"}
}

type resourceHandler struct {
	pathTemplate string
	dir          string
}

func (h *resourceHandler) Register(r *mux.Router) {
	r.HandleFunc(h.pathTemplate, h.handle)
}

func (h *resourceHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := writeResource(h.dir, w, r); err != nil {
		http.NotFound(w, r)
		return
	}
}

func writeResource(resourceDir string, w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("writeResource(): %v", urlPath)

	relativePath := mux.Vars(r)["path"]
	resource, err := resources.Resolve(resourceDir, relativePath)
	if err != nil {
		return err
	}

	// If this fails, we've already started writing the response, so it's too
	// late to return a 404 or whatever; just log it and move on
	err = resource.Write(w, urlPath)
	log.Print(err)
	return nil
}
