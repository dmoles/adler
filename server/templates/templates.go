package templates

import (
	"fmt"
	resources2 "github.com/dmoles/adler/resources"
	"log"
	"path"
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
	Header string
	Title  string
	TOC    string
	Body   string
}

// ------------------------------------------------------------
// Unexported

func load(name string) *template.Template {
	tmplPath := path.Join("/templates", name)
	resource, err := resources2.Get(tmplPath)
	if err != nil {
		msg := fmt.Sprintf("Error locating template %s: %v", tmplPath, err)
		log.Fatal(msg)
	}

	tmplData, err := resource.AsString()
	if err != nil {
		msg := fmt.Sprintf("Error reading template %s: %v", tmplPath, err)
		log.Fatal(msg)
	}

	tmpl, err := template.New(name).Parse(tmplData)
	if err != nil {
		msg := fmt.Sprintf("Error parsing template %s: %v", tmplPath, err)
		log.Fatal(msg)
	}
	return tmpl
}
