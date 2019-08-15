package adler

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"io/ioutil"
	"os"
	"path"
)

// Embeddable sub-fixture for temporary directories
type tempdirFixture struct {
	*gunit.Fixture
	tempDir string
}

func (f *tempdirFixture) SetupTempDir() {
	tempDir, err := ioutil.TempDir("", "tempdirFixture")
	f.So(err, should.BeNil)
	f.tempDir = tempDir
}

func (f *tempdirFixture) TeardownTempDir() {
	err := os.RemoveAll(f.tempDir)
	f.So(err, should.BeNil)
}

// Create a subdirectory relative to the fixture's temporary
// directory, and return its path as a string
func (f *tempdirFixture) CreateDir(relativeDir string) string {
	dirActual := path.Join(f.tempDir, relativeDir)
	err := os.MkdirAll(dirActual, 0755)
	f.So(err, should.BeNil)
	return dirActual
}

// Create an empty file at a path relative to the fixture's temporary
// directory, and return its path as a string
func (f *tempdirFixture) CreateFile(relativePath string) string {
	relativeDir := path.Dir(relativePath)
	dirActual := f.CreateDir(relativeDir)

	baseName := path.Base(relativePath)
	pathActual := path.Join(dirActual, baseName)

	file, err := os.Create(pathActual)
	f.So(err, should.BeNil)
	err = file.Close()
	f.So(err, should.BeNil)

	return pathActual
}