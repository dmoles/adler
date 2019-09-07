package adler

import (
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
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
	port    int
}

func (s *server) Start() error {
	log.Printf("Serving from %s on port %d", s.RootDir(), s.port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	var addr = fmt.Sprintf(":%d", s.port)
	return graceful.RunWithErr(addr, finishRequestTimeout, mux)

}

// TODO: figure out why relative links don't work
func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	// TODO: generate TOC sidebar

	mdPath, err := s.Resolve(r.URL.Path[1:])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	page, err := NewPage(mdPath)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	data := pageData { page.Title(), page.ToHtml() }
	err = pageTemplate.Execute(w, data)

	if err != nil {
		log.Println(err)
	}
}

type pageData struct {
	Title string
	Body string
}

var pageTemplate = func() *template.Template {
	headTmpl := trim(`
	<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>
	<main>
	{{.Body}}
	</main>
	</body>
	</html>
	`)

	tmpl, err := template.New("page").Parse(headTmpl)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}()

//// Deprecated: TODO: reimplement w/Resolver
//func (s *server) handler(w http.ResponseWriter, r *http.Request) {
//	path := r.URL.Path[1:]
//
//	var pages []string
//	err := filepath.Walk(s.rootDir, func(path string, info os.FileInfo, err error) error {
//		// TODO: handle subdirectories properly
//		name := info.Name()
//		if strings.HasSuffix(name, ".md") {
//			pages = append(pages, strings.TrimSuffix(name, ".md"))
//		}
//		return nil
//	})
//
//	// TODO: generate directory indexes
//
//	// TODO: handle URL paths that already end in `.md`
//
//	mdPath := filepath.Join(s.rootDir, path+".md")
//	if _, err := os.Stat(mdPath); err != nil {
//		if os.IsNotExist(err) {
//			http.Error(w, err.Error(), http.StatusNotFound)
//			return
//		}
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	content, err := ioutil.ReadFile(mdPath)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	output := blackfriday.Run(content)
//
//	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
//
//	// TODO: use a template language
//	// TODO: some CSS
//	_, _ = fmt.Fprintln(w, "<html>")
//	_, _ = fmt.Fprintln(w, "<head>")
//	_, _ = fmt.Fprintf(w, "  <title>%s</title>\n", path)
//	_, _ = fmt.Fprintln(w, "</head>")
//	_, _ = fmt.Fprintln(w, "<body>")
//	_, _ = fmt.Fprintln(w, "<main>")
//
//	_, _ = w.Write(output)
//
//	_, _ = fmt.Fprintln(w, "</main>")
//
//	_, _ = fmt.Fprintln(w, "<aside>")
//	_, _ = fmt.Fprintln(w, "<ul>")
//	for _, page := range pages {
//		_, _ = fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>\n", page, page)
//	}
//	_, _ = fmt.Fprintln(w, "</ul>")
//	_, _ = fmt.Fprintln(w, "</aside>")
//
//	_, _ = fmt.Fprintln(w, "</body>")
//	_, _ = fmt.Fprintln(w, "</html>")
//}
