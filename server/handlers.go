package server

import (
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/resources"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

func resourceHandler(dir string, varname string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		relativePath := vars[varname]
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

// TODO: should this just be a function on Server?
func rawHandler(rootDir string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		log.Printf("raw(): %v", urlPath)
		filePath, err := util.ResolveFile(urlPath, rootDir)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		// TODO: stat first then stream
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		util.WriteData(w, urlPath, data)
	}
}

// TODO: should this just be a function on Server?
func markdownHandler(rootDir string) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path

		resolvedPath, err := util.ResolvePath(urlPath, rootDir)
		if err != nil {
			log.Printf("Error resolving path %v: %v", urlPath, err)
			http.NotFound(w, r)
			return
		}

		title, err := markdown.GetTitleFromFile(resolvedPath)
		if err != nil {
			log.Printf("Error determining title from path: %v: %v", resolvedPath, err)
			http.NotFound(w, r)
			return
		}

		rootIndexHtml, err := markdown.DirToHtml(rootDir, rootDir)
		if err != nil {
			log.Printf("Error generating directory index for %v: %v", rootDir, err)
			http.NotFound(w, r)
			return
		}

		bodyHtml, err := markdown.GetBodyHTML(resolvedPath, rootDir)

		pageData := templates.PageData{
			Title: title,
			TOC:   string(rootIndexHtml),
			Body:  string(bodyHtml),
		}

		var sb strings.Builder
		err = templates.Page().Execute(&sb, pageData)
		if err != nil {
			log.Printf("Error executing template for %v: %v", urlPath, err)
			http.NotFound(w, r)
			return
		}

		data := []byte(sb.String())
		util.WriteData(w, urlPath, data)

	}
}
