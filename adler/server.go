package adler

import (
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const finishRequestTimeout = 5 * time.Second

type Server interface {
	Start() error
}

func NewServer(port int, rootDir string) (Server, error) {
	resolver, err := NewResolver(rootDir)
	if err != nil {
		return nil, err
	}
	return &server{resolver, port}, nil
}

type server struct {
	Resolver
	port int
}

func (s *server) Start() error {
	log.Printf("Serving from %s on port %d", s.RootDir(), s.port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	var addr = fmt.Sprintf(":%d", s.port)
	return graceful.RunWithErr(addr, finishRequestTimeout, mux)

}

// TODO: handle images etc.
func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	// TODO: generate TOC sidebar

	path := r.URL.Path
	if s.isCSS(path) {
		err := s.serveCSS(w, path)
		if err != nil {
			s.error404(w, err)
		}
		return
	}

	err := s.serveHTML(w, path)
	if err != nil {
		s.error404(w, err)
	}
}

func (s *server) logError(err error) {
	log.Println(err)
}

func (s *server) error404(w http.ResponseWriter, err error) {
	s.logError(err)
	http.Error(w, err.Error(), http.StatusNotFound)
}

func (s *server) isCSS(path string) bool {
	return strings.HasPrefix(path, "/css/")
}

func (s *server) serveCSS(w http.ResponseWriter, path string) error {
	cssPath := strings.TrimPrefix(path, "/css")
	data, err := findCSS(cssPath)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "text/css; charset=UTF-8")
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return err
	}
	return nil
}

func (s *server) serveHTML(w http.ResponseWriter, path string) error {
	mdPath, err := s.Resolve(path[1:])
	if err != nil {
		return err
	}

	// TODO: something smarter
	//       - title
	//       - subdirectories
	tocPage, err := NewPage(s.RootDir())
	if err != nil {
		return err
	}

	page, err := NewPage(mdPath)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	// TODO: make sure all possible errors happen before first body write
	data := pageData{
		TOC: tocPage.ToHtml(),
		Title: page.Title(),
		Body: page.ToHtml(),
	}
	return pageTemplate.Execute(w, data)
}

type pageData struct {
	TOC string
	Title string
	Body  string
}

// TODO: move this to a file cf. CSS
var pageTemplate = func() *template.Template {
	headTmpl := trim(`
	<html>
	<head>
		<title>{{.Title}}</title>
		<link rel="stylesheet" href="/css/reset.css">
		<link rel="stylesheet" href="/css/main.css">
	</head>
	<header>
	<nav>
		<h5><a href="/">Home</a></h5>
	</nav>
	</header>
	<body>
	<main>
	{{.Body}}
	</main>
	<aside>
	<nav>
		{{.TOC}}
	</nav>
	</aside>
	</body>
	</html>
	`)

	tmpl, err := template.New("page").Parse(headTmpl)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}()
