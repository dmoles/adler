package server

import (
	"fmt"
	"github.com/dmoles/adler/server/handlers"
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
		return err;
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

func (s *server) newRouter() *mux.Router {
	// TODO: support single-page version
	r := mux.NewRouter()
	// TODO: look up handlers automatically from path pattern/prefix
	r.PathPrefix("/css/{css:.+}").HandlerFunc(s.css)
	r.HandleFunc("/{favicon:[^/]+\\.(?:ico|png|jpg)}", s.favicon)
	r.HandleFunc("/{markdown:.+\\.md}", s.markdown)
	r.MatcherFunc(s.isDirectory).HandlerFunc(s.markdown)
	r.MatcherFunc(s.isFile).HandlerFunc(s.raw)
	return r
}

// ------------------------------
// Handler methods

func (s *server) css(w http.ResponseWriter, r *http.Request) {
	// TODO: use SCSS
	handlers.CSS().Handle(w, r)
}

func (s *server) favicon(w http.ResponseWriter, r *http.Request) {
	handlers.Favicon().Handle(w, r)
}

func (s *server) markdown(w http.ResponseWriter, r *http.Request) {
	handlers.Markdown(s.rootDir).Handle(w, r)
}

func (s *server) raw(w http.ResponseWriter, r *http.Request) {
	log.Printf("raw(): %v", r.URL.Path)
	handlers.Raw(s.rootDir).Handle(w, r)
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
