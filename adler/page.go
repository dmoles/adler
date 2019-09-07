package adler

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Page interface {
	Title() string
	Content() []byte
	RelativeLink() string
	ToHtml() string
}

func NewPage(filePath string) (Page, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return newIndexPage(filePath)
	}
	return newMarkdownPage(filePath)
}

// ------------------------------------------------------------
// page

type page struct {
	title    string
	content []byte
	filePath string
	html string
}

func newPage(path string, content []byte) *page {
	title := textOfFirstHeading(content)
	if title == "" {
		title = asTitle(path)
	}
	return &page{title: title, content: content, filePath: path}
}

func (p *page) Title() string {
	return p.title
}

func (p *page) Content() []byte {
	return p.content
}

func (p *page) RelativeLink() string {
	return fmt.Sprintf("[%v](%v)", p.title, path.Base(p.filePath))
}

func (p *page) ToHtml() string {
	if p.html == "" {
		p.html = string(blackfriday.Run(p.content))
	}
	return p.html
}

// ------------------------------------------------------------
// Initializers

func newMarkdownPage(filePath string) (*page, error) {
	if !isMarkdownFile(filePath) {
		return nil, fmt.Errorf("%#v is not a Markdown file", filePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return newPage(filePath, content), nil
}


func newIndexPage(dirPath string) (*page, error) {
	if !isDirectory(dirPath) {
		return nil, fmt.Errorf("%#v is not a directory", dirPath)
	}

	var sb strings.Builder
	_, err := fmt.Fprintf(&sb, "# %s\n\n", asTitle(dirPath))
	if err != nil {
		return nil, err
	}

	var links []string
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, info := range files {
		fullPath := filepath.Join(dirPath, info.Name())
		if page, err := NewPage(fullPath); err == nil {
			link := page.RelativeLink()
			links = append(links, link)
		}
	}
	sort.Strings(links)
	for _, link := range links {
		_, err = fmt.Fprintf(&sb, "- %v\n", link)
		if err != nil {
			return nil, err
		}
	}

	return newPage(dirPath, []byte(sb.String())), nil
}