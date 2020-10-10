package templates

import (
	"fmt"
	"github.com/dmoles/adler/server/resources"
	"io"
	"log"
	"path"
	"strings"
	"text/template"
)

// ------------------------------------------------------------
// Exported

const (
	page = "page.html.tmpl"
)

var pageTemplate *template.Template

func Page() *template.Template {
	if pageTemplate == nil {
		pageTemplate = load(page)
	}
	return pageTemplate
}

type PageData struct {
	Title string
	TOC   string
	Body  string
}

// ------------------------------------------------------------
// Unexported

func load(name string) *template.Template {
	tmplPath := path.Join("/templates", name)
	tmplFile, err := resources.Open(tmplPath)
	if err != nil {
		msg := fmt.Sprintf("Error locating template %s: %v", tmplPath, err)
		log.Fatal(msg)
	}
	sb := new(strings.Builder)
	_, err = io.Copy(sb, tmplFile)
	if err != nil {
		msg := fmt.Sprintf("Error reading template %s: %v", tmplPath, err)
		log.Fatal(msg)
	}
	tmplData := sb.String()
	tmpl, err := template.New(name).Parse(tmplData)
	if err != nil {
		msg := fmt.Sprintf("Error parsing template %s: %v", tmplPath, err)
		log.Fatal(msg)
	}
	return tmpl
}
