package handlers

import (
	"github.com/bep/golibsass/libsass"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	cssPathTemplate = "/css/{path:.+\\.css}"
	cssResourceDir  = "/css"
)

func CSS(cssDir string) Handler {
	if cssDir == "" {
		return &resourceHandler{cssPathTemplate, cssResourceDir}
	}
	return &localCSSHandler{cssDir: cssDir}
}

type localCSSHandler struct {
	cssDir string
	transpilerP libsass.Transpiler
}

func (h *localCSSHandler) Register(r *mux.Router) {
	r.HandleFunc(cssPathTemplate, h.handle)
}

func (h *localCSSHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := h.serveCss(w, r); err != nil {
		log.Printf("no CSS or SCSS file found for URL path %v: %v", r.URL.Path, err)
		http.NotFound(w, r)
	}
}

func (h *localCSSHandler) serveCss(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("serveCss(): %v", urlPath)
	relativePath := mux.Vars(r)["path"]

	filePath, err := util.UrlPathToFile(relativePath, h.cssDir)
	if err == nil {
		return writeRaw(filePath, w, r)
	}
	scssUrlPath := strings.TrimSuffix(relativePath, ".css") + ".scss"
	scssFilePath, err := util.UrlPathToFile(scssUrlPath, h.cssDir)
	if err == nil {
		return h.serveScss(scssFilePath, w, r)
	}
	return err
}

func (h *localCSSHandler) transpiler() (libsass.Transpiler, error) {
	if h.transpilerP == nil {
		libsassOptions := libsass.Options{
			IncludePaths: []string{h.cssDir},
			OutputStyle:  libsass.ExpandedStyle,
		}
		transpilerP, err := libsass.New(libsassOptions)
		if err != nil {
			return nil, err
		}
		h.transpilerP = transpilerP
	}
	return h.transpilerP, nil
}

func (h *localCSSHandler) serveScss(scssFilePath string, w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	log.Printf("serveScss(#%v): %v", scssFilePath, urlPath)

	buf, err := ioutil.ReadFile(scssFilePath)
	if err != nil {
		return err
	}
	scssSrc := string(buf)

	transpiler, err := h.transpiler()
	if err != nil {
		return err
	}

	result, err := transpiler.Execute(scssSrc)
	if err != nil {
		return err
	}
	cssData := []byte(result.CSS)

	// If this fails, we've already started writing the response, so it's too
	// late to return a 404 or whatever; just log it and move on
	err = util.WriteData(w, urlPath, cssData)
	log.Print(err)
	return nil
}
