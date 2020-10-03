package handlers

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type markdownHandler struct {
	rootDir string
}

func (h *markdownHandler) Handle(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	rootDir := h.rootDir

	resolvedPath, err := util.ResolvePath(urlPath, rootDir)
	if err != nil {
		log.Printf("Error resolving path %v: %v", urlPath, err)
		http.NotFound(w, r)
		return
	}

	title, err := markdown.GetTitleFromFile(resolvedPath)
	if err != nil {
		log.Printf("Error determining title from path: %v: %v", resolvedPath, err)
		http.NotFound(w, r)
		return
	}

	rootIndexHtml, err := markdown.DirToHtml(h.rootDir, h.rootDir)
	if err != nil {
		log.Printf("Error generating directory index for %v: %v", rootDir, err)
		http.NotFound(w, r)
		return
	}

	bodyHtml, err := h.GetBodyHtml(resolvedPath)

	pageData := templates.PageData{
		Title: title,
		TOC:   string(rootIndexHtml),
		Body:  string(bodyHtml),
	}

	var sb strings.Builder
	err = templates.Page().Execute(&sb, pageData)
	if err != nil {
		log.Printf("Error executing template for %v: %v", urlPath, err)
		http.NotFound(w, r)
		return
	}

	data := []byte(sb.String())
	util.WriteData(w, urlPath, data)
}

func (h *markdownHandler) GetBodyHtml(resolvedPath string) ([]byte, error) {
	if util.IsDirectory(resolvedPath) {
		readmePath := filepath.Join(resolvedPath, "README.md")
		if util.IsFile(readmePath) {
			resolvedPath = readmePath
		} else {
			return markdown.DirToHtml(resolvedPath, h.rootDir)
		}
	}
	return markdown.FileToHtml(resolvedPath), nil
}
