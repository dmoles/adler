package templates

import (
	"github.com/markbates/pkger"
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
	TOC string
	Body  string
}

// ------------------------------------------------------------
// Unexported

func load(name string) *template.Template {
	tmplPath := path.Join("/resources/templates", name)
	tmplFile, err := pkger.Open(tmplPath)
	if err != nil {
		log.Fatal(err)
	}
	sb := new(strings.Builder)
	_, err = io.Copy(sb, tmplFile)
	if err != nil {
		log.Fatal(err)
	}
	tmplData := sb.String()
	tmpl, err := template.New(name).Parse(tmplData)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}
