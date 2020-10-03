package handlers

import (
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// TODO: centralize resource utility code & also use for templates
// TODO: in makefile: pkger -include /templates -include /css -include /images

type resourceHandler struct {
	dir string;
	varname string;
}

func (h *resourceHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	relativePath := vars[h.varname]
	if relativePath == "" || strings.Contains(relativePath, "..") {
		http.NotFound(w, r)
		return
	}
	relativePathClean := path.Clean(relativePath)
	resourcePath := path.Join(h.dir, relativePathClean)
	file, err := pkger.Open(resourcePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentType := util.ContentType(resourcePath)
	w.Header().Add("Content-Type", contentType)

	size := fileInfo.Size()
	contentLength := strconv.FormatInt(size, 10)
	w.Header().Add("Content-Length", contentLength)

	n, err := io.Copy(w, file)
	if err != nil {
		log.Printf("Error serving %#v: %v", resourcePath, err)
	}
	if n != size {
		log.Printf("Wrote wrong number of bytes for %#v: expected %d, was %d", resourcePath, size, n)
	}
}
