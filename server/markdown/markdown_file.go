package markdown

import (
	"github.com/dmoles/adler/server/markdown"
)

type MarkdownFile struct {
	urlPath string
	title string
	bodyHtml []byte
	stylesheets []string
	scripts []Script
}

func MarkdownFileFrom(resolvedPath string) (*MarkdownFile, error) {
	bodyHtml, metadata, err := markdown.FileToHtml(resolvedPath)
	if err != nil {
		return nil, err
	}

	mf := &MarkdownFile {
		bodyHtml: bodyHtml,
	}
	return mf, nil
}
