package adler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Deprecated: TODO: test & integrate w/Resolver
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
	_, err = fmt.Fprintf(&sb, "# %s\n\n", p.Title())
	if err != nil {
		return nil, err
	}
	for _, info := range files {
		if link, ok := relativeLink(p.dirPath, info); ok {
			_, err = fmt.Fprintf(&sb, "- %v\n", link)
			if err != nil {
				return nil, err
			}
		}
	}
	return []byte(sb.String()), nil
}


