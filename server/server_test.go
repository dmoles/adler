package server

import (
	"github.com/dmoles/adler/server/util"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

const (
	invalidPort = -1
)

func testdataDir() string {
	projectRoot := util.ProjectRoot()
	return filepath.Join(projectRoot, "testdata")
}

func newServer(t *testing.T) Server {
	s, err := New(invalidPort, testdataDir(), "")
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return s
}

func get(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return req
}

func TestMarkdown(t *testing.T) {
	expect := NewWithT(t).Expect
	recorder := httptest.NewRecorder()
	router := newServer(t).(*server).newRouter()

	router.ServeHTTP(recorder, get(t, "/hello.md"))

	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/html; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("<title>Hello</title>"))
	expect(body).To(ContainSubstring("<h1>Hello</h1>"))
	expect(body).To(ContainSubstring("<p>Hello, world</p>"))
}

func TestCSSResource(t *testing.T) {
	expect := NewWithT(t).Expect
	recorder := httptest.NewRecorder()
	router := newServer(t).(*server).newRouter()

	router.ServeHTTP(recorder, get(t, "/css/main.css"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("text/css; charset=utf-8"))

	body := recorder.Body.String()
	expect(body).To(HavePrefix("@charset \"UTF-8\";"))
}

func TestFavicon(t *testing.T) {
	expect := NewWithT(t).Expect
	recorder := httptest.NewRecorder()
	router := newServer(t).(*server).newRouter()

	router.ServeHTTP(recorder, get(t, "/favicon.ico"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("image/x-icon"))
}

func TestFooterIcon(t *testing.T) {
	expect := NewWithT(t).Expect
	recorder := httptest.NewRecorder()
	router := newServer(t).(*server).newRouter()

	router.ServeHTTP(recorder, get(t, "/apple-touch-icon.png"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("image/png"))
}

func TestWOFFFont(t *testing.T) {
	expect := NewWithT(t).Expect
	recorder := httptest.NewRecorder()
	router := newServer(t).(*server).newRouter()

	router.ServeHTTP(recorder, get(t, "/fonts/CharisSIL-R.woff"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("font/woff"))
}

func TestWOFF2Font(t *testing.T) {
	expect := NewWithT(t).Expect
	recorder := httptest.NewRecorder()
	router := newServer(t).(*server).newRouter()

	router.ServeHTTP(recorder, get(t, "/fonts/CourierPrime-Regular.woff2"))
	expect(recorder.Code).To(Equal(http.StatusOK))

	result := recorder.Result()
	contentTypes := result.Header["Content-Type"]
	expect(contentTypes).To(HaveLen(1))
	expect(contentTypes[0]).To(Equal("font/woff2"))
}
