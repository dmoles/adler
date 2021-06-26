package handlers

import (
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	cssPathTemplate = "/css/{path:.+\\.css}"
	cssResourceDir  = "/css"
)

func CSS(cssDir string) Handler {
	if cssDir == "" {
		return &resourceHandler{cssPathTemplate, cssResourceDir}
	}
	return &localCSSHandler{cssDir: cssDir}
}

type localCSSHandler struct {
	cssDir string
}

func (h *localCSSHandler) Register(r *mux.Router) {
	r.HandleFunc(cssPathTemplate, h.handle)
}

func (h *localCSSHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := h.serveCss(w, r); err != nil {
		log.Printf("can't serve CSS for URL path %v: %v", r.URL.Path, err)
		http.NotFound(w, r)
	}
}

func (h *localCSSHandler) serveCss(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("serveCss(): %v", urlPath)
	relativePath := mux.Vars(r)["path"]

	filePath, err := util.UrlPathToFile(relativePath, h.cssDir)
	if err == nil {
		return writeRaw(filePath, w, r)
	}
	return err
}
