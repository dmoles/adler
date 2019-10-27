package templates

import (
	"github.com/gobuffalo/packr"
	"log"
	"text/template"
)

// ------------------------------------------------------------
// Exported

const (
	page = "page.html.tmpl"
)

var Page = load(page)

type PageData struct {
	Title string
	TOC string
	Body  string
}

// ------------------------------------------------------------
// Unexported

var templateBox = packr.NewBox("../../templates")

func load(name string) *template.Template {
	tmplData, err := templateBox.Find(name)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New(name).Parse(string(tmplData))
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}
