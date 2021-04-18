package server

import (
	"context"
	"fmt"
	"github.com/dmoles/adler/server/handlers"
	"github.com/dmoles/adler/server/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// ------------------------------------------------------------
// Exported

type Server interface {
	Start() error
}

func (s *server) Start() error {
	log.Printf("Serving from %s on port %d", s.rootDir, s.port)
	if s.cssDir != "" {
		log.Printf("Using CSS directory %v", s.cssDir)
	}
	router := s.newRouter()

	addr := fmt.Sprintf(":%d", s.port)

  srv := &http.Server {
		Handler: router,
		Addr: addr,
		WriteTimeout: finishRequestTimeout,
		ReadTimeout: finishRequestTimeout,
	}
	
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func Start(port int, rootDir string, cssDir string) error {
	s, err := New(port, rootDir, cssDir)
	if err != nil {
		return err
	}
	return s.Start()
}

func New(port int, rootDir string, cssDir string) (Server, error) {
	return newServer(port, rootDir, cssDir)
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

func (s *server) newRouter() *mux.Router {
	// TODO: support single-page version
	r := mux.NewRouter()
	for _, h := range handlers.All(s.rootDir, s.cssDir) {
		h.Register(r)
	}
	return r
}

// ------------------------------
// Unexported functions

func newServer(port int, rootDir string, cssDir string) (*server, error) {
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
