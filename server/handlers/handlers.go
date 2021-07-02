package handlers

import (
	"github.com/gorilla/mux"
)

// ------------------------------
// Exported symbols

type Handler interface {
	Register(r *mux.Router)
}

func All(rootDir string) []Handler {
	return []Handler{
		MarkdownFile(rootDir),
		DirectoryIndex(rootDir),
		Raw(rootDir),
		ResourceHandler(),
	}
}
