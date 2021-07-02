package handlers

import (
	"fmt"
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)


const markdownPathPattern = "/{path:.+\\.md}"

func MarkdownFile(rootDir string) Handler {
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

// TODO: something less awful
func (h *markdownHandler) writeFile(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	//log.Printf("write(): %v", urlPath)

	resolvedPath, err := util.UrlPathToFile(urlPath, h.rootDir)
	if err != nil {
		return err
	}

	bodyHtml, metadata, err := markdown.FileToHtml(resolvedPath)
	if err != nil {
		return err
	}

	var title string

	titleVal := metadata["Title"]
	if titleVal == nil {
		title, err = markdown.ExtractTitle(resolvedPath)
		if err != nil {
			return err
		}
	} else {
		var ok bool
		title, ok = titleVal.(string)
		if !ok {
			return fmt.Errorf("not a string: %#v", titleVal)
		}
	}

	var stylesheets []string
	stylesheetsVal := metadata["Stylesheets"]
	if stylesheetsVal != nil {
		stylesheetVals, ok := stylesheetsVal.([]interface{})
		if !ok {
			return fmt.Errorf("not a []interface{}: %#v", stylesheetsVal)
		}
		for _, stylesheetVal := range stylesheetVals {
			stylesheet, ok := stylesheetVal.(string)
			if !ok {
				return fmt.Errorf("not a string: %#v", stylesheet)
			}
			stylesheets = append(stylesheets, stylesheet)
		}
	}

	var scripts []string
	scriptsVal := metadata["Scripts"]
	if scriptsVal != nil {
		scriptVals, ok := scriptsVal.([]interface{})
		if !ok {
			return fmt.Errorf("not a []interface{}: %#v", scriptsVal)
		}
		for _, scriptVal := range scriptVals {
			script, ok := scriptVal.(string)
			if !ok {
				return fmt.Errorf("not a string: %#v", script)
			}
			scripts = append(scripts, script)
		}
	}

	return h.write(w, urlPath, title, stylesheets, scripts, bodyHtml)
}

type markdownHandlerBase struct {
	rootDir string
}

func (h *markdownHandlerBase) write(w http.ResponseWriter, urlPath string, title string, stylesheets []string, scripts []string, bodyHtml []byte) error {
	rootIndexHtml, _, err := markdown.DirToIndexHtml(h.rootDir, h.rootDir)
	if err != nil {
		return err
	}

	siteTitle, err := markdown.GetTitleFromFile(h.rootDir)
	if err != nil {
		return err
	}

	pageData := templates.PageData{
		Header: siteTitle,
		Title:  title,
		TOC:    string(rootIndexHtml),
		Body:   string(bodyHtml),
		Stylesheets: stylesheets,
		Scripts: scripts,
	}

	var sb strings.Builder
	err = templates.Page().Execute(&sb, pageData)
	if err != nil {
		return err
	}
	return util.WriteData(w, urlPath, []byte(sb.String()))
}
