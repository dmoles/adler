package adler

import (
	"fmt"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/lithammer/dedent"
)

// ------------------------------------------------------------
// Page

func TestPageFixture(t *testing.T) {
	gunit.Run(new(PageFixture), t)
}

type PageFixture struct {
	*gunit.Fixture
	TempdirManager
}

func (f *PageFixture) Setup() {
	err := f.InitTempDir()
	f.So(err, should.BeNil)
}

func (f *PageFixture) TestMarkdownPage() {
	body := dedent.Dedent(`
		# Expected title
	
		Body of document
	`)

	filePath, err := f.CreateFile("path/to/file.md")
	f.So(err, should.BeNil)
	err = ioutil.WriteFile(filePath, []byte(body), 0644)
	f.So(err, should.BeNil)

	page, err := NewPage(filePath)
	f.So(err, should.BeNil)

	f.So(page.Title(), should.Equal, "Expected title")

	content, err := page.Content()
	f.So(err, should.BeNil)
	f.So(string(content), should.Equal, body)
}

func (f *PageFixture) TestIndexPage() {
	dir, err := f.CreateDir("dirName")
	f.So(err, should.BeNil)

	for i := 1; i <= 2; i++ {
		body := []byte(fmt.Sprintf("# File %d\n", i))
		baseName := fmt.Sprintf("file%d.md", i)
		filePath := filepath.Join(dir, baseName)
		err = ioutil.WriteFile(filePath, body, 0644)
		f.So(err, should.BeNil)
	}

	expectedBody := trim(`
	# DirName

	- [File 1](file1.md)
	- [File 2](file2.md)
	`)

	page, err := NewPage(dir)
	f.So(err, should.BeNil)

	f.So(page.Title(), should.Equal, "DirName")

	content, err := page.Content()
	f.So(err, should.BeNil)
	f.So(string(content), should.Equal, expectedBody)
}

func (f *PageFixture) TestIndexPageSupportsSubdirectories() {
	dir, err := f.CreateDir("dirName")
	f.So(err, should.BeNil)

	for i := 1; i <= 2; i++ {
		body := []byte(fmt.Sprintf("# File %d\n", i))
		baseName := fmt.Sprintf("file%d.md", i)
		filePath := filepath.Join(dir, baseName)
		err = ioutil.WriteFile(filePath, body, 0644)
		f.So(err, should.BeNil)
	}

	_, err = f.CreateDir("dirName/subdirectory")
	f.So(err, should.BeNil)

	expectedBody := trim(`
	# DirName

	- [File 1](file1.md)
	- [File 2](file2.md)
	- [Subdirectory](subdirectory)
	`)

	page, err := NewPage(dir)
	f.So(err, should.BeNil)

	f.So(page.Title(), should.Equal, "DirName")

	content, err := page.Content()
	f.So(err, should.BeNil)
	f.So(string(content), should.Equal, expectedBody)
}
