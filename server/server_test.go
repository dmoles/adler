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

func TestRouter(t *testing.T) {
	expect := NewWithT(t).Expect
	
	s, err := New(invalidPort, testdataDir())
	expect(err).NotTo(HaveOccurred())

	router := s.(*server).newRouter()

	req, err := http.NewRequest("GET", "/hello.md", nil)
	expect(err).NotTo(HaveOccurred())

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	expect(recorder.Code).To(Equal(http.StatusOK))

	body := recorder.Body.String()
	expect(body).To(ContainSubstring("Hello, world"))
}
