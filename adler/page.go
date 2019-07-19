package adler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Page interface {
	Title() string
	Content() ([]byte, error)
}

func NewPage(filePath string) (Page, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return newIndexPage(filePath), nil
	}
	return newMarkdownPage(filePath)
}

// ------------------------------------------------------------
// markdownPage

type markdownPage struct {
	title    string
	filePath string
}

func newMarkdownPage(filePath string) (*markdownPage, error) {
	if !isMarkdownFile(filePath) {
		return nil, fmt.Errorf("%#v is neither a Markdown file nor a directory", filePath)
	}

	title, err := textOfFirstHeading(filePath)
	if err != nil {
		return nil, err
	}
	if title == "" {
		title = filepath.Base(filePath)
	}
	return &markdownPage{title: title, filePath: filePath}, nil
}

func (p *markdownPage) Title() string {
	return p.title
}

func (p *markdownPage) Content() ([]byte, error) {
	return ioutil.ReadFile(p.filePath)
}

// ------------------------------------------------------------
// indexPage

type indexPage struct {
	dirPath string
}

func newIndexPage(dirPath string) *indexPage {
	return &indexPage{dirPath: dirPath}
}

func (p *indexPage) Title() string {
	return filepath.Base(p.dirPath)
}

func (p *indexPage) Content() ([]byte, error) {
	var sb strings.Builder
	files, err := ioutil.ReadDir(p.dirPath)
	if err != nil {
		return nil, err
	}
	for _, info := range files {
		if link, ok := relativeLink(p.dirPath, info); ok {
			//noinspection GoUnhandledErrorResult
			fmt.Fprintf(&sb, "- %v\n", link)
		}
	}
	return []byte(sb.String()), nil
}


