package adler

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Page interface {
	Title() string
	Content() []byte
	RelativeLink() string
	ToHtml() string
}

func NewPageFromPath(filePath string) (Page, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return newIndexPage(filePath)
	}
	return newMarkdownPage(filePath)
}

func NewPage(title string, content []byte, basePath string) Page {
	return &page{title: title, content: content, basePath: basePath}
}

// ------------------------------------------------------------
// page

type page struct {
	title    string
	content  []byte
	basePath string
	html     string
}

func newPage(realPath string, content []byte) Page {
	title := textOfFirstHeading(content)
	if title == "" {
		title = asTitle(realPath)
	}
	basePath := path.Base(realPath)
	return NewPage(title, content, basePath)
}

func (p *page) Title() string {
	return p.title
}

func (p *page) Content() []byte {
	return p.content
}

func (p *page) RelativeLink() string {
	return fmt.Sprintf("[%v](%v)", p.title, p.basePath)
}

func (p *page) ToHtml() string {
	if p.html == "" {
		// TODO: figure out AutoHeadingIds
		p.html = string(blackfriday.Run(p.content))
	}
	return p.html
}

type singlePage struct {
	title    string
	pages []Page
	basePath string
}

func (p *singlePage) Title() string {
	return p.title
}

func (p *singlePage) Content() []byte {
	// TODO straighten out these types
	panic("not implemented")
}

func (p *singlePage) RelativeLink() string {
	return fmt.Sprintf("[%v](%v)", p.title, p.basePath)
}

func (p *singlePage) ToHtml() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "<h1>%s</h1>\n\n", p.title)
	for i, p := range(p.pages) {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString("<article>\n")
		sb.WriteString(p.ToHtml())
		sb.WriteString("</article>\n")
	}
	return sb.String()
}

// ------------------------------------------------------------
// Initializers

func newMarkdownPage(filePath string) (Page, error) {
	if !isMarkdownFile(filePath) {
		return nil, fmt.Errorf("%#v is not a Markdown file", filePath)
	}


	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return newPage(filePath, content), nil
}


func newIndexPage(dirPath string) (Page, error) {
	if !isDirectory(dirPath) {
		return nil, fmt.Errorf("%#v is not a directory", dirPath)
	}

	var sb strings.Builder
	_, err := fmt.Fprintf(&sb, "# %s\n\n", asTitle(dirPath))
	if err != nil {
		return nil, err
	}

	var titles []string
	var links = map[string]string{}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, info := range files {
		filename := info.Name()
		if isHidden(filename) {
			continue
		}
		fullPath := filepath.Join(dirPath, filename)
		if page, err := NewPageFromPath(fullPath); err == nil {
			link := page.RelativeLink()
			title := page.Title()
			titles = append(titles, title)
			links[title] = link
		}
	}
	sort.Slice(titles, func(i, j int) bool {
		st1 := sortingTitle(titles[i])
		st2 := sortingTitle(titles[j])
		return st1 < st2;
	})
	for _, title := range titles {
		link := links[title]
		_, err = fmt.Fprintf(&sb, "- %v\n", link)
		if err != nil {
			return nil, err
		}
	}

	return newPage(dirPath, []byte(sb.String())), nil
}

var numericPrefixRegexp = regexp.MustCompile("^[0-9-]+ (.+)")

func sortingTitle(t string) string {
	st := strings.TrimSpace(strings.ToLower(t));

	if submatch := numericPrefixRegexp.FindStringSubmatch(st); submatch != nil {
		return submatch[1]
	}

	for _, prefix := range []string{"a ", "the "} {
		if strings.HasPrefix(st, prefix) {
			return strings.TrimPrefix(st, prefix)
		}
	}

	return st;
}

func NewSinglePage(dirPath string) (Page, error) {
	if !isDirectory(dirPath) {
		return nil, fmt.Errorf("%#v is not a directory", dirPath)
	}

	title := asTitle(dirPath)

	var pages []Page
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, info := range files {
		filename := info.Name()
		if isHidden(filename) {
			log.Printf("skipping hidden file %v\n", filename)
			continue
		}
		fullPath := filepath.Join(dirPath, filename)
		if !isMarkdownFile(fullPath) {
			log.Printf("skipping non-Markdown file %v\n", filename)
			continue
		}
		page, err := NewPageFromPath(fullPath);
		if err == nil {
			log.Printf("adding Markdown file %v\n", filename)
			pages = append(pages, page)
		} else {
			log.Printf("error adding file %v: %v\n", filename, err)
		}
	}

	sort.Slice(pages, func(i, j int) bool {
		st1 := sortingTitle(pages[i].Title())
		st2 := sortingTitle(pages[j].Title())
		return st1 < st2;
	})
	return &singlePage{title, pages, path.Base(dirPath)}, nil
}
