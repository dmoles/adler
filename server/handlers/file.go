package handlers

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type fileHandler struct {
	rootDir string
	serve func(w http.ResponseWriter, r *http.Request, path string)
}

func (h *fileHandler) handle(w http.ResponseWriter, r *http.Request) {
	rootDir := h.rootDir
	urlPath := r.URL.Path
	resolvedPath, err := util.ResolvePath(urlPath, rootDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.serve(w, r, resolvedPath)
}

func serveRaw(w http.ResponseWriter, r *http.Request, path string) {
	urlPath := r.URL.Path
	log.Printf("raw(): %v", urlPath)

	filePath, err := util.ToAbsoluteFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	util.WriteData(w, urlPath, data)
}

func serveMarkdown(w http.ResponseWriter, r *http.Request, path string) {
	urlPath := r.URL.Path
	resolvedPath := path

	title, err := markdown.GetTitleFromFile(resolvedPath)
	if err != nil {
		log.Printf("Error determining title from path: %v: %v", resolvedPath, err)
		http.NotFound(w, r)
		return
	}

	rootIndexHtml, err := markdown.DirToHtml(rootDir, rootDir)
	if err != nil {
		log.Printf("Error generating directory index for %v: %v", rootDir, err)
		http.NotFound(w, r)
		return
	}

	bodyHtml, err := markdown.GetBodyHTML(resolvedPath, rootDir)

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