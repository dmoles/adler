package server

import (
	"github.com/dmoles/adler/server/util"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	. "github.com/onsi/gomega"
)

const (
	invalidPort = -1
)

func testdataDir() string {
	projectRoot := util.ProjectRoot()
	return filepath.Join(projectRoot, "testdata")
}

func TestRouter(t *testing.T) {
	s, err := New(invalidPort, testdataDir())
	Expect(err).NotTo(HaveOccurred())

	router := s.(*server).newRouter()

	req, err := http.NewRequest("GET", "/hello.md", nil)
	Expect(err).NotTo(HaveOccurred())

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	Expect(recorder.Code).To(Equal(http.StatusOK))

	body := recorder.Body.String()
	Expect(body).To(ContainSubstring("Hello, world"))
}
