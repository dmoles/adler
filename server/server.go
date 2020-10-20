package server

import (
	"fmt"
	"github.com/dmoles/adler/server/markdown"
	"github.com/dmoles/adler/server/templates"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server interface {
	Start() error
}

func Start(port int, rootDir string) error {
	s, err := New(port, rootDir)
	if err != nil {
		return err
	}
	return s.Start()
}

func New(port int, rootDir string) (Server, error) {
	rootDirAbs, err := util.ToAbsoluteDirectory(rootDir)
	if err != nil {
		return nil, err
	}
	return &server{port: port, rootDir: rootDirAbs}, nil
}

// ------------------------------------------------------------
// Unexported

const finishRequestTimeout = 5 * time.Second

type server struct {
	port    int
	rootDir string
	router  *mux.Router
}

func (s *server) Start() error {
	log.Printf("Serving from %s on port %d", s.rootDir, s.port)
	router := s.newRouter()

	addr := fmt.Sprintf(":%d", s.port)
	return graceful.RunWithErr(addr, finishRequestTimeout, router)
}

// ------------------------------
// Private functions

const markdownPathPattern = "/{path:.+\\.md}"

func (s *server) newRouter() *mux.Router {
	// TODO: support single-page version
	r := mux.NewRouter()

	for pathTemplate, handler := range resourceHandlers {
		r.HandleFunc(pathTemplate, handler)
	}

	r.HandleFunc(markdownPathPattern, s.handleMarkdown)
	r.MatcherFunc(s.isDirectory).HandlerFunc(s.handleMarkdown)

	r.MatcherFunc(s.isFile).HandlerFunc(s.handleRaw)
	return r
}

// ------------------------------
// Handlers

func (s *server) handleRaw(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	log.Printf("raw(): %v", urlPath)
	filePath, err := util.ResolveFile(urlPath, s.rootDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// TODO: just stream, we already checked for existence
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	util.WriteData(w, urlPath, data)
}

func (s *server) handleMarkdown(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	resolvedPath, err := util.ResolvePath(urlPath, s.rootDir)
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

	rootIndexHtml, err := markdown.DirToHtml(s.rootDir, s.rootDir)
	if err != nil {
		log.Printf("Error generating directory index for %v: %v", s.rootDir, err)
		http.NotFound(w, r)
		return
	}

	bodyHtml, err := markdown.GetBodyHTML(resolvedPath, s.rootDir)

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

// ------------------------------
// Utility methods

func (s *server) inRootDir(urlPath string) bool {
	_, err := util.ResolveFile(urlPath, s.rootDir)
	return err == nil
}

func (s *server) isDirectory(r *http.Request, rm *mux.RouteMatch) bool {
	_, err := util.ResolveDirectory(r.URL.Path, s.rootDir)
	return err == nil
}

func (s *server) isFile(r *http.Request, rm *mux.RouteMatch) bool {
	_, err := util.ResolveFile(r.URL.Path, s.rootDir)
	return err == nil
}
