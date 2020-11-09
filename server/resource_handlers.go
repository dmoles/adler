package server

import (
	"github.com/dmoles/adler/server/resources"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

var resourceHandlers = map[string]handlerFunc{
	"/css/{path:.+}":   makeHandler("/css"),
	"/fonts/{path:.+}": makeHandler("/fonts"),
	"/{path:[^/]+\\.(?:ico|png|jpg|webmanifest)}": makeHandler("/images/favicons"),
}

func makeHandler(dir string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		relativePath := vars["path"]
		if relativePath == "" || strings.Contains(relativePath, "..") {
			http.NotFound(w, r)
			return
		}
		relativePathClean := path.Clean(relativePath)
		resourcePath := path.Join(dir, relativePathClean)

		resource, err := resources.Get(resourcePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		contentType := resource.ContentType()
		w.Header().Add("Content-Type", contentType)

		size := resource.Size()
		w.Header().Add("Content-Length", strconv.FormatInt(size, 10))

		n, err := resource.Copy(w)
		if err != nil {
			log.Printf("Error serving %#v: %v", resourcePath, err)
		}
		if n != size {
			log.Printf("Wrote wrong number of bytes for %#v: expected %d, was %d", resourcePath, size, n)
		}
	}
}
