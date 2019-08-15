package adler

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"testing"
)

func TestResolver(t *testing.T) {
	gunit.Run(new(NewResolverFixture), t)
}

type NewResolverFixture struct {
	*gunit.Fixture
	tempDir string
}

func (f *NewResolverFixture) Setup() {
	tempDir, err := ioutil.TempDir("", "NewResolverFixture")
	f.So(err, should.BeNil)
	f.tempDir = tempDir
}

func (f *NewResolverFixture) Teardown() {
	err := os.RemoveAll(f.tempDir)
	f.So(err, should.BeNil)
}

func (f *NewResolverFixture) TestSetsAbsoluteRootDir() {
	rootDir := path.Join(f.tempDir, "root")
	rootDirAbs, err := filepath.Abs(rootDir)
	f.So(err, should.BeNil)
	err = os.Mkdir(rootDirAbs, 0755)
	f.So(err, should.BeNil)

	r, err := NewResolver(rootDir)
	f.So(err, should.BeNil)
	f.So(r, should.NotBeNil)
	f.So(r.rootDir, should.Equal, rootDirAbs)
}

func (f *NewResolverFixture) TestRootDirMustExist() {
	rootDir := path.Join(f.tempDir, "root")
	err := os.RemoveAll(rootDir)
	f.So(err, should.BeNil)

	r, err := NewResolver(rootDir)
	f.So(r, should.BeNil)
	f.So(err, should.NotBeNil)
}

func (f *NewResolverFixture) TestRootDirMustBeADirectory() {
	rootDir := path.Join(f.tempDir, "root")
	file, err := os.Create(rootDir)
	f.So(err, should.BeNil)
	err = file.Close()
	f.So(err, should.BeNil)

	r, err := NewResolver(rootDir)
	f.So(r, should.BeNil)
	f.So(err, should.NotBeNil)
}
