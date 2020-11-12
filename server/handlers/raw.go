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
	err := writeRaw(h.resolvePath, w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *rawHandler) resolvePath(r *http.Request) (string, error) {
	urlPath := r.URL.Path
	resolvedPath, err := util.ResolveRelative(urlPath, h.rootDir)
	if err != nil {
		return "", err
	}

	return util.ToAbsoluteFile(resolvedPath)
}

type pathResolver func(r *http.Request) (string, error)

func writeRaw(resolvePath pathResolver, w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("writeRaw(): %v", urlPath)

	filePath, err := resolvePath(r)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// If this fails, we've already started writing the response, so it's too
	// late to return a 404 or whatever; just log it and move on
	err = util.WriteData(w, urlPath, data)
	log.Print(err)
	return nil
}
