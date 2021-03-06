package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"

	"github.com/dmoles/adler/server/util"
)

const (
	invalidPort = -1
)

// ------------------------------
// Fixture

type expectFunc func(actual interface{}, extra ...interface{}) Assertion

func testdataDir() string {
	projectRoot := util.ProjectRoot()
	return filepath.Join(projectRoot, "testdata")
}

func get(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return req
}

func setUp(t *testing.T) (expect expectFunc, recorder *httptest.ResponseRecorder, router *mux.Router) {
	return setUpWithCSS(t, "")
}

func setUpWithCSS(t *testing.T, cssDir string) (expect expectFunc, recorder *httptest.ResponseRecorder, router *mux.Router) {
	serverP, err := newServer(invalidPort, testdataDir(), cssDir)
	if err != nil {
		t.Fatal(err)
		return
	}

	router = serverP.newRouter()
	expect = NewWithT(t).Expect
	recorder = httptest.NewRecorder()

	return
}

// ------------------------------
// Tests

func TestMarkdown(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/hello.md"))

	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/html; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>Hello</title>"))
	expect(body).To(ContainSubstring("<h1 id=\"hello\">Hello</h1>"))
	expect(body).To(ContainSubstring("<p>Hello, world</p>"))
}

func TestMarkdownMetadata(t *testing.T) {
	expect, recorder, router := setUp(t)

	expectedTitle := "The Real Title"
	basename := "metadata-test.md"

	urlPath := path.Join("/", basename)
	router.ServeHTTP(recorder, get(t, urlPath))
	expect(recorder.Code).To(Equal(http.StatusOK))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>" + expectedTitle + "</title>"))

	expectedStylesheets := []string{"foo.css", "css/bar.css"}
	for _, stylesheet := range expectedStylesheets {
		expectedLink := fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"/>", stylesheet)
		expect(body).To(ContainSubstring(expectedLink))
	}

	expectedScripts := []string{
		"<script src=\"foo.js\" type=\"module\"></script>",
		"<script src=\"scripts/bar.js\"></script>",
	}
	for _, script := range expectedScripts {
		expect(body).To(ContainSubstring(script))
	}
}

func TestDirectoryIndex(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/html; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>Testdata</title>"))

	navRe := regexp.MustCompile("(?s)<nav>.*</nav>")
	nav := navRe.FindString(body)
	expect(nav).NotTo(BeEmpty())
	expect(nav).To(ContainSubstring("<li><a href=\"/hello.md\">Hello</a></li>"))

	mainRe := regexp.MustCompile("(?s)<main>.*</main>")
	main := mainRe.FindString(body)
	expect(main).NotTo(BeEmpty())
	expect(main).To(ContainSubstring("<li><a href=\"/hello.md\">Hello</a></li>"))
}

func TestSubdirectoryIndex(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/subdirectory/"))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/html; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>Subdirectory</title>"))

	// nav should still be there in subdirectories
	navRe := regexp.MustCompile("(?s)<nav>.*</nav>")
	nav := navRe.FindString(body)
	expect(nav).NotTo(BeEmpty())
	expect(nav).To(ContainSubstring("<li><a href=\"/hello.md\">Hello</a></li>"))

	mainRe := regexp.MustCompile("(?s)<main>.*</main>")
	main := mainRe.FindString(body)
	expect(main).NotTo(BeEmpty())
	expect(main).To(ContainSubstring("<li><a href=\"/subdirectory/hello-again.md\">hello again</a></li>"))
}

func TestSubdirectoryWithoutTrailingSlash(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/subdirectory/"))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/html; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>Subdirectory</title>"))

	// nav should still be there in subdirectories
	navRe := regexp.MustCompile("(?s)<nav>.*</nav>")
	nav := navRe.FindString(body)
	expect(nav).NotTo(BeEmpty())
	expect(nav).To(ContainSubstring("<li><a href=\"/hello.md\">Hello</a></li>"))

	mainRe := regexp.MustCompile("(?s)<main>.*</main>")
	main := mainRe.FindString(body)
	expect(main).NotTo(BeEmpty())
	expect(main).To(ContainSubstring("<li><a href=\"/subdirectory/hello-again.md\">hello again</a></li>"))
}

func TestReadme(t *testing.T) {
	expect, recorder, router := setUp(t)

	readmeData := "# Testing\n\nTesting 123"
	readmePath := filepath.Join(testdataDir(), "README.md")
	err := ioutil.WriteFile(readmePath, []byte(readmeData), 0600)
	defer util.RemoveAllQuietly(readmePath)
	expect(err).NotTo(HaveOccurred()) // just to be sure

	router.ServeHTTP(recorder, get(t, "/"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/html; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>Testing</title>"))

	navRe := regexp.MustCompile("(?s)<nav>.*</nav>")
	nav := navRe.FindString(body)
	expect(nav).NotTo(BeEmpty())
	expect(nav).To(ContainSubstring("<li><a href=\"/hello.md\">Hello</a></li>"))

	mainRe := regexp.MustCompile("(?s)<main>.*</main>")
	main := mainRe.FindString(body)
	expect(main).NotTo(BeEmpty())
	expect(main).To(ContainSubstring("<h1 id=\"testing\">Testing</h1>"))
	expect(main).To(ContainSubstring("<p>Testing 123</p>"))
	expect(main).NotTo(ContainSubstring("<li><a href=\"/README.md\">Testing</a></li>"))
}

func TestCSSResource(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/css/main.css"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/css; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(MatchRegexp("^@charset \"UTF-8\";\\s+html, body, div, span"))
}

func TestCSSOverride(t *testing.T) {
	// Setup

	cssDir := filepath.Join(testdataDir(), "css")
	err := os.Mkdir(cssDir, 0700)
	if err == nil {
		defer util.RemoveAllQuietly(cssDir)
	} else {
		t.Error(err)
	}

	cssData := "body { background-color: #808000; }"
	cssPath := filepath.Join(cssDir, "main.css")

	err = ioutil.WriteFile(cssPath, []byte(cssData), 0600)
	if err != nil {
		t.Error(err)
	}

	expect, recorder, router := setUpWithCSS(t, cssDir)

	// Test serving local CSS
	router.ServeHTTP(recorder, get(t, "/css/main.css"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/css; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(Equal(cssData))
}

func TestFavicon(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/favicon.ico"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("image/x-icon"))
}

func TestFooterIcon(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/apple-touch-icon.png"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("image/png"))
}

func TestWOFFFont(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/fonts/CharisSIL-R.woff"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("font/woff"))
}

func TestWOFF2Font(t *testing.T) {
	expect, recorder, router := setUp(t)

	router.ServeHTTP(recorder, get(t, "/fonts/CourierPrime-Regular.woff2"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("font/woff2"))
}
