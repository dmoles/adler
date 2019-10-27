package handlers

import (
	"github.com/gobuffalo/packr"
	"net/http"
)

// ------------------------------------------------------------
// Exported

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func CSS() Handler {
	return &boxHandler{box: cssBox, varname: "css"}
}

func Favicon() Handler {
	return &boxHandler{box: faviconBox, varname: "favicon"}
}

func Raw(rootDir string) Handler {
	return &rawHandler{rootDir: rootDir}
}

func Markdown(rootDir string) Handler {
	return &markdownHandler{rootDir: rootDir}
}

// ------------------------------------------------------------
// Unexported

var cssBox = packr.NewBox("../../css")
var faviconBox = packr.NewBox("../../images/favicons")
