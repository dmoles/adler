package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/util"
)

const markdownPathPattern = "/{path:.+\\.md}"

func Markdown(rootDir string) Handler {
	h := markdownHandler{}
	h.rootDir = rootDir
	return &h
}

type markdownHandler struct {
	markdownHandlerBase
}

func (h *markdownHandler) Register(r *mux.Router) {
	r.HandleFunc(markdownPathPattern, h.handle)
}

func (h *markdownHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.writeFile(w, r)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
	}
}

func (h *markdownHandler) writeFile(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	//log.Printf("write(): %v", urlPath)

	resolvedPath, err := util.UrlPathToFile(urlPath, h.rootDir)
	if err != nil {
		return err
	}

	mf, err := markdown.FromFile(resolvedPath)
	if err != nil {
		return err
	}

	return h.write(w, urlPath, mf)
}
