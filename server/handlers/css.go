package handlers

import (
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	cssPathTemplate = "/css/{path:.+}"
	cssResourceDir  = "/css"
)

func CSS(cssDir string) Handler {
	if cssDir == "" {
		return &resourceHandler{cssPathTemplate, cssResourceDir}
	}
	return &cssRawHandler{cssDir: cssDir}
}

type cssRawHandler struct {
	cssDir string
}

func (h *cssRawHandler) Register(r *mux.Router) {
	r.HandleFunc(cssPathTemplate, h.handle)
}

func (h *cssRawHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := writeRaw(h.resolvePath, w, r); err != nil {
		http.NotFound(w, r)
		return
	}
}

func (h *cssRawHandler) resolvePath(r *http.Request) (string, error) {
	relativePath := mux.Vars(r)["path"]
	resolvedPath, err := util.ResolveRelative(relativePath, h.cssDir)
	if err != nil {
		return "", err
	}

	return util.ToAbsoluteFile(resolvedPath)
}
