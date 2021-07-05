package handlers

import (
	"net/http"
	"strings"

	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
)

type markdownHandlerBase struct {
	rootDir string
}

func (h *markdownHandlerBase) write(w http.ResponseWriter, urlPath string, mf markdown.MarkdownFile) error {
	resolvedPath, err := util.ResolveUrlPath(urlPath, h.rootDir)
	rootIndex, err := markdown.DirectoryIndex(h.rootDir, resolvedPath)
	if err != nil {
		return err
	}

	rootIndexHtml := rootIndex.MainContent().ToHtml()
	siteTitle := rootIndex.Title().Text()

	var headElements []string
	for _, h := range mf.Head() {
		headElements = append(headElements, h.ToHtml())
	}

	pageData := templates.PageData{
		Header:       siteTitle,
		Title:        mf.Title().Text(),
		HeadElements: headElements,
		TOC:          rootIndexHtml,
		Body:         mf.MainContent().ToHtml(),
	}

	var sb strings.Builder
	err = templates.Page().Execute(&sb, pageData)
	if err != nil {
		return err
	}

	return util.WriteData(w, urlPath, []byte(sb.String()))
}
