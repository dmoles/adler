package handlers

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	_, err := util.ResolveDirectory(r.URL.Path, h.rootDir)
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
	log.Printf("write(): %v", urlPath)

	rootDir := h.rootDir

	resolvedPath, err := util.ResolveDirectory(urlPath, rootDir)
	if err != nil {
		return err
	}

	title := markdown.AsTitle(resolvedPath)

	bodyHtml, err := markdown.DirToHTML(resolvedPath, rootDir)
	if err != nil {
		return err
	}

	return h.write(w, urlPath, title, bodyHtml)
}
