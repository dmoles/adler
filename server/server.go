package server

import (
	"fmt"
	"github.com/dmoles/adler/server/handlers"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"path/filepath"
	"time"
)

type Server interface {
	Start() error
}

func Start(port int, rootDir string, cssDir string) error {
	s, err := New(port, rootDir, cssDir)
	if err != nil {
		return err
	}
	return s.Start()
}

func New(port int, rootDir string, cssDir string) (Server, error) {
	rootDirAbs, err := util.ToAbsoluteDirectory(rootDir)
	if err != nil {
		return nil, err
	}

	newServer := &server{port: port, rootDir: rootDirAbs}

	if cssDir != "" {
		cssDirAbs, err := util.ToAbsoluteDirectory(cssDir)
		if err != nil {
			return nil, err
		}
		main_css := filepath.Join(cssDirAbs, "main.css")
		main_scss := filepath.Join(cssDirAbs, "main.scss")
		if util.IsFile(main_css) != util.IsFile(main_scss) {
			newServer.cssDir = cssDirAbs
		} else {
			return nil, fmt.Errorf("CSS directory %v must contain exactly one of: main.css, main.scss", cssDirAbs)
		}
	}

	return newServer, nil
}

// ------------------------------------------------------------
// Unexported

const finishRequestTimeout = 5 * time.Second

type server struct {
	port    int
	rootDir string
	router  *mux.Router
	cssDir  string
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
	for _, h := range handlers.All(s.rootDir) {
		h.Register(r)
	}
	return r
}
