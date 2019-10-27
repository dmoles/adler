package handlers

import (
	"github.com/dmoles/adler/server/util"
	"io/ioutil"
	"net/http"
)

type rawHandler struct {
	rootDir string
}

func (h *rawHandler) Handle(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	path, err := util.ResolveFile(urlPath, h.rootDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// TODO: stream this
	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	util.WriteData(w, urlPath, data)
}