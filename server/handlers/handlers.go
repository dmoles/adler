package handlers

import (
	"net/http"
)

// ------------------------------------------------------------
// Exported

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func CSS() Handler {
	return &resourceHandler{dir: "/resources/css", varname: "css"}
}

func Favicon() Handler {
	return &resourceHandler{dir: "/resources/images/favicons", varname: "favicon"}
}

func Raw(rootDir string) Handler {
	return &rawHandler{rootDir: rootDir}
}

func Markdown(rootDir string) Handler {
	return &markdownHandler{rootDir: rootDir}
}
