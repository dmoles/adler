package handlers

import (
	"github.com/gorilla/mux"
)

// ------------------------------
// Exported symbols

type Handler interface {
	Register(r *mux.Router)
}

func All(rootDir string, cssDir string) []Handler {
	return []Handler{
		CSS(cssDir),
		FontResource(),
		FaviconResource(),
		MarkdownFile(rootDir),
		DirectoryIndex(rootDir),
		Raw(rootDir),
	}
}
