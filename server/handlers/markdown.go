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

const markdownPathPattern = "/{path:.+\\.md}"

type markdownHandler struct {
	rootDir string
}

func (h *markdownHandler) Register(r *mux.Router) {
	r.HandleFunc(markdownPathPattern, h.handle)
}

func (h *markdownHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.writeMarkdown(w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

// TODO: DRY with directoryHandler
func (h *markdownHandler) writeMarkdown(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("writeMarkdown(): %v", urlPath)

	rootDir := h.rootDir

	resolvedPath, err := util.ResolveFile(urlPath, rootDir)
	if err != nil {
		return err
	}

	title, err := markdown.GetTitleFromFile(resolvedPath)
	if err != nil {
		return err
	}

	rootIndexHtml, err := markdown.DirToIndexHtml(rootDir, rootDir)
	if err != nil {
		return err
	}

	bodyHtml, err := markdown.FileToHtml(resolvedPath)
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