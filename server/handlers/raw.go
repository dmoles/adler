package handlers

import (
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func Raw(rootDir string) Handler {
	// TODO: simplify this -- see https://github.com/gorilla/mux#static-files
	return &rawHandler{rootDir: rootDir}
}

type rawHandler struct {
	rootDir string
}

func (h *rawHandler) Register(r *mux.Router) {
	r.MatcherFunc(h.isFile).HandlerFunc(h.handle)
}

func (h *rawHandler) isFile(r *http.Request, _ *mux.RouteMatch) bool {
	_, err := util.UrlPathToFile(r.URL.Path, h.rootDir)
	return err == nil
}

func (h *rawHandler) handle(w http.ResponseWriter, r *http.Request) {
	err := h.serveRaw(w, r)
	if err != nil {
		http.NotFound(w, r)
	}
}

func (h *rawHandler) serveRaw(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	//log.Printf("serveRaw(): %v", urlPath)

	filePath, err := util.UrlPathToFile(urlPath, h.rootDir)
	if err != nil {
		return err
	}

	return writeRaw(filePath, w, r)
}

func writeRaw(filePath string, w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	//log.Printf("writeRaw(%#v): %v", filePath, urlPath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// If this fails, we've already started writing the response, so it's too
	// late to return a 404 or whatever; just log it and move on
	err = util.WriteData(w, urlPath, data)
	if err != nil {
		log.Print(err)
	}
	return nil
}
