package handlers

import (
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func Raw(rootDir string) Handler {
	return &rawHandler{rootDir: rootDir}
}

type rawHandler struct {
	rootDir string
}

func (h *rawHandler) Register(r *mux.Router) {
	r.MatcherFunc(h.isFile).HandlerFunc(h.handle)
}

func (h *rawHandler) isFile(r *http.Request, _ *mux.RouteMatch) bool {
	_, err := util.ResolveFile(r.URL.Path, h.rootDir)
	return err == nil
}

func (h *rawHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.writeRaw(w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *rawHandler) writeRaw(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("writeRaw(): %v", urlPath)

	resolvedPath, err := util.ResolveRelative(urlPath, h.rootDir)
	if err != nil {
		return err
	}

	filePath, err := util.ToAbsoluteFile(resolvedPath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	util.WriteData(w, urlPath, data)

	return nil
}
