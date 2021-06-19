package handlers

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"net/http"
	"strings"
)

type markdownHandlerBase struct {
	rootDir string
}

func (h *markdownHandlerBase) write(w http.ResponseWriter, urlPath string, title string, bodyHtml []byte) error {
	rootIndexHtml, err := markdown.DirToIndexHtml(h.rootDir, h.rootDir)
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
	}

	var sb strings.Builder
	err = templates.Page().Execute(&sb, pageData)
	if err != nil {
		return err
	}
	return util.WriteData(w, urlPath, []byte(sb.String()))
}
