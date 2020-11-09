package handlers

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const markdownPathPattern = "/{path:.+\\.md}"

func MarkdownFile(rootDir string) Handler {
	h := fileHandler{}
	h.rootDir = rootDir
	return &h
}

type fileHandler struct {
	markdownHandlerBase
}

func (h *fileHandler) Register(r *mux.Router) {
	r.HandleFunc(markdownPathPattern, h.handle)
}

func (h *fileHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.writeFile(w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *fileHandler) writeFile(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("write(): %v", urlPath)

	resolvedPath, err := util.ResolveFile(urlPath, h.rootDir)
	if err != nil {
		return err
	}

	title, err := markdown.ExtractTitle(resolvedPath)
	if err != nil {
		return err
	}

	bodyHtml, err := markdown.FileToHtml(resolvedPath)
	if err != nil {
		return err
	}

	return h.write(w, urlPath, title, bodyHtml)
}
