package adler

import (
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
	"path"
	"path/filepath"
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

	urlPath := r.URL.Path
	if s.isCSS(urlPath) {
		err := s.serveCSS(w, urlPath)
		if err != nil {
			s.error404(w, err)
		}
		return
	}

	if s.isFavicon(urlPath) {
		err := s.serveFavicon(w, urlPath)
		if err != nil {
			s.error404(w, err)
		}
		return
	}

	if s.isSinglePage(urlPath) {
		err := s.serveSinglePage(w, urlPath)
		if err != nil {
			s.error404(w, err)
		}
		return
	}

	if isMarkdownFile(urlPath) || isDirectory(urlPath) {
		err := s.serveHTML(w, urlPath)
		if err != nil {
			s.error404(w, err)
		}
		return
	}

	err := s.serveImage(w, urlPath)
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

func (s *server) isCSS(urlPath string) bool {
	return strings.HasPrefix(urlPath, "/css/")
}

func (s *server) isSinglePage(urlPath string) bool {
	return strings.HasPrefix(urlPath, "/single")
}

func (s *server) serveCSS(w http.ResponseWriter, urlPath string) error {
	log.Printf("serveCSS(%#v)", urlPath)

	cssPath := strings.TrimPrefix(urlPath, "/css")
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

func (s *server) isFavicon(urlPath string) bool {
	if "/favicon.ico" == urlPath {
		return true
	}
	dir := path.Dir(urlPath)
	if dir != "/" {
		return false;
	}
	if !strings.HasSuffix(urlPath, ".png") {
		return false
	}
	if !strings.Contains(urlPath, "icon") {
		return false
	}
	return true
}

func (s *server) serveFavicon(w http.ResponseWriter, urlPath string) error {
	log.Printf("serveFavicon(%#v)", urlPath)

	imagePath := path.Join("favicons/", path.Base(urlPath))
	data, contentType, err := findImage(imagePath)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", contentType)
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return err
	}
	return nil
}

func (s *server) serveSinglePage(w http.ResponseWriter, path string) error {
	log.Printf("serveSinglePage(%#v)", path)

	dirPath, err := s.Resolve(path)
	if err != nil {
		return err
	}
	dirPath = filepath.Dir(dirPath)
	page, err := NewSinglePage(dirPath)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	data := pageData {
		TOC: "",
		Title: page.Title(),
		Body: page.ToHtml(),
	}
	return singlePageTemplate.Execute(w, data)
}

func (s *server) serveImage(w http.ResponseWriter, path string) error {
	log.Printf("serveImage(%#v)", path)

	rawPath, err := s.Resolve(path)
	if err != nil {
		return err
	}
	data, contentType, err := findImage(rawPath)
	if err != nil {
		return err;
	}
	w.Header().Add("Content-Type", contentType)
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
	log.Printf("serveHTML(%#v)", path)
	
	mdPath, err := s.Resolve(path[1:])
	if err != nil {
		return err
	}

	// TODO: something smarter
	//       - title
	//       - subdirectories
	tocPage, err := NewPageFromPath(s.RootDir())
	if err != nil {
		return err
	}

	page, err := NewPageFromPath(mdPath)
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

		<link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png">
		<link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png">
		<link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png">
		<link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png">
		<link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png">
		<link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png">
		<link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png">
		<link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png">
		<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
		<link rel="icon" type="image/png" sizes="192x192"  href="/android-icon-192x192.png">
		<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
		<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
		<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
		<link rel="manifest" href="/manifest.json">
		<meta name="msapplication-TileColor" content="#ffffff">
		<meta name="msapplication-TileImage" content="/ms-icon-144x144.png">
		<meta name="theme-color" content="#ffffff">
	</head>
	<header>
	<nav>
		<h5><a href="/">Home</a></h5>
	</nav>
	</header>
	<body>
	<aside>
	<nav>
		{{.TOC}}
	</nav>
	</aside>
	<main>
		{{.Body}}
	</main>
	<footer>
	<p><img class="adler-icon" src="/favicon-96x96.png"/> Served by <a href="https://github.com/dmoles/adler/">Adler</a>.</p>
	</footer>
	</body>
	</html>
	`)

	tmpl, err := template.New("page").Parse(headTmpl)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}()

var singlePageTemplate = func() *template.Template {
	headTmpl := trim(`
	<html>
	<head>
		<title>{{.Title}}</title>
		<link rel="stylesheet" href="/css/reset.css">
		<link rel="stylesheet" href="/css/main.css">

		<link rel="apple-touch-icon" sizes="57x57" href="/apple-icon-57x57.png">
		<link rel="apple-touch-icon" sizes="60x60" href="/apple-icon-60x60.png">
		<link rel="apple-touch-icon" sizes="72x72" href="/apple-icon-72x72.png">
		<link rel="apple-touch-icon" sizes="76x76" href="/apple-icon-76x76.png">
		<link rel="apple-touch-icon" sizes="114x114" href="/apple-icon-114x114.png">
		<link rel="apple-touch-icon" sizes="120x120" href="/apple-icon-120x120.png">
		<link rel="apple-touch-icon" sizes="144x144" href="/apple-icon-144x144.png">
		<link rel="apple-touch-icon" sizes="152x152" href="/apple-icon-152x152.png">
		<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
		<link rel="icon" type="image/png" sizes="192x192"  href="/android-icon-192x192.png">
		<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
		<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
		<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
		<link rel="manifest" href="/manifest.json">
		<meta name="msapplication-TileColor" content="#ffffff">
		<meta name="msapplication-TileImage" content="/ms-icon-144x144.png">
		<meta name="theme-color" content="#ffffff">
	</head>
	<body>
	<main class="full">
	{{.Body}}
	</main>
	<footer>
	<p><img class="adler-icon" src="/favicon-96x96.png"/> Served by <a href="https://github.com/dmoles/adler/">Adler</a>.</p>
	</footer>
	</body>
	</html>
	`)

	tmpl, err := template.New("page").Parse(headTmpl)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}()

