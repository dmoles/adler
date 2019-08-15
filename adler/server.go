package adler

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const finishRequestTimeout = 5 * time.Second

type Server interface {
	Start() error
}

func NewServer(port int, rootDir string) (Server, error) {
	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}
	return &server{port: port, rootDir: rootDirAbs}, nil
}

type server struct {
	port    int
	rootDir string
}

func (s *server) Start() error {
	log.Printf("Serving from %s on port %d", s.rootDir, s.port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler)

	var addr = fmt.Sprintf(":%d", s.port)
	return graceful.RunWithErr(addr, finishRequestTimeout, mux)
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	var pages []string
	err := filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
		// TODO: handle subdirectories properly
		name := info.Name()
		if strings.HasSuffix(name, ".md") {
			pages = append(pages, strings.TrimSuffix(name, ".md"))
		}
		return nil
	})

	// TODO: generate directory indexes

	// TODO: handle URL paths that already end in `.md`

	mdPath := filepath.Join(s.rootDir, path+".md")
	if _, err := os.Stat(mdPath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	output := blackfriday.Run(content)

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	// TODO: use a template language
	// TODO: some CSS
	_, _ = fmt.Fprintln(w, "<html>")
	_, _ = fmt.Fprintln(w, "<head>")
	_, _ = fmt.Fprintf(w, "  <title>%s</title>\n", path)
	_, _ = fmt.Fprintln(w, "</head>")
	_, _ = fmt.Fprintln(w, "<body>")
	_, _ = fmt.Fprintln(w, "<main>")

	_, _ = w.Write(output)

	_, _ = fmt.Fprintln(w, "</main>")

	_, _ = fmt.Fprintln(w, "<aside>")
	_, _ = fmt.Fprintln(w, "<ul>")
	for _, page := range pages {
		_, _ = fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>\n", page, page)
	}
	_, _ = fmt.Fprintln(w, "</ul>")
	_, _ = fmt.Fprintln(w, "</aside>")

	_, _ = fmt.Fprintln(w, "</body>")
	_, _ = fmt.Fprintln(w, "</html>")
}
