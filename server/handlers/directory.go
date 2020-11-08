package handlers

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type directoryHandler struct {
	rootDir string
}

func (h *directoryHandler) Register(r *mux.Router) {
	r.MatcherFunc(h.isDirectory).HandlerFunc(h.handle)
}

func (h *directoryHandler) isDirectory(r *http.Request, _ *mux.RouteMatch) bool {
	// TODO: DRY util.ResolveDirectory
	_, err := util.ResolveDirectory(r.URL.Path, h.rootDir)
	return err == nil
}

func (h *directoryHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.writeDirectory(w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

// TODO: DRY with markdownHandler
func (h *directoryHandler) writeDirectory(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("writeDirectory(): %v", urlPath)

	rootDir := h.rootDir

	resolvedPath, err := util.ResolveDirectory(urlPath, rootDir)
	if err != nil {
		return err
	}

	title := markdown.AsTitle(resolvedPath)

	rootIndexHtml, err := markdown.DirToIndexHtml(rootDir, rootDir)
	if err != nil {
		return err
	}

	bodyHtml, err := markdown.DirToHTML(resolvedPath, rootDir)
	if err != nil {
		return err
	}

	pageData := templates.PageData{
		Title: title,
		TOC:   string(rootIndexHtml),
		Body:  string(bodyHtml),
	}

	var sb strings.Builder
	err = templates.Page().Execute(&sb, pageData)
	if err != nil {
		return nil
	}
	util.WriteData(w, urlPath, []byte(sb.String()))

	return nil
}