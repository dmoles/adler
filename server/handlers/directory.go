package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/util"
)

func DirectoryIndex(rootDir string) Handler {
	h := directoryHandler{}
	h.rootDir = rootDir
	return &h
}

type directoryHandler struct {
	markdownHandlerBase
}

func (h *directoryHandler) Register(r *mux.Router) {
	r.MatcherFunc(h.isDirectory).HandlerFunc(h.handle)
}

func (h *directoryHandler) isDirectory(r *http.Request, _ *mux.RouteMatch) bool {
	_, err := util.UrlPathToDirectory(r.URL.Path, h.rootDir)
	return err == nil
}

func (h *directoryHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.writeDirectory(w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *directoryHandler) writeDirectory(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	//log.Printf("write(): %v", urlPath)

	rootDir := h.rootDir

	resolvedPath, err := util.UrlPathToDirectory(urlPath, rootDir)
	if err != nil {
		return err
	}

	mf, err := markdown.ForDirectory(resolvedPath, rootDir)
	if err != nil {
		return err
	}

	return h.write(w, urlPath, mf)
}
