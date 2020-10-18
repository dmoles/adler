package server

import (
	"fmt"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
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

const cssPathPrefix = "/css/{css:.+}"
const fontsPathPrefix = "/fonts/{fonts:.+}"
const faviconPathPattern = "/{favicon:[^/]+\\.(?:ico|png|jpg)}"
const markdownPathPattern = "/{markdown:.+\\.md}"

func (s *server) newRouter() *mux.Router {
	// TODO: support single-page version
	r := mux.NewRouter()

	// TODO: can we unify on either PathPrefix or HandleFunc & DRY?
	r.PathPrefix(cssPathPrefix).HandlerFunc(resourceHandler("/css", "css"))
	r.PathPrefix(fontsPathPrefix).HandlerFunc(resourceHandler("/fonts", "fonts"))
	r.HandleFunc(faviconPathPattern, resourceHandler("/images/favicons", "favicon"))

	markdown := markdownHandler(s.rootDir)
	r.HandleFunc(markdownPathPattern, markdown)
	r.MatcherFunc(s.isDirectory).HandlerFunc(markdown)

	raw := rawHandler(s.rootDir)
	r.MatcherFunc(s.isFile).HandlerFunc(raw)
	return r
}

// ------------------------------
// Utility methods

func (s *server) isDirectory(r *http.Request, rm *mux.RouteMatch) bool {
	_, err := util.ResolveDirectory(r.URL.Path, s.rootDir)
	return err == nil
}

func (s *server) isFile(r *http.Request, rm *mux.RouteMatch) bool {
	_, err := util.ResolveFile(r.URL.Path, s.rootDir)
	return err == nil
}
