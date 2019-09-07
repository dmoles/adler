package adler

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"os"
	"path"
	"path/filepath"

	"testing"
)

// ------------------------------------------------------------
// Resolver

func TestResolverFixture(t *testing.T) {
	gunit.Run(new(ResolverFixture), t)
}

type ResolverFixture struct {
	*gunit.Fixture
	TempdirManager
	resolver Resolver
}

func (f *ResolverFixture) Setup() {
	err := f.InitTempDir()
	f.So(err, should.BeNil)

	resolver, err := NewResolver(f.tempDir)
	f.So(err, should.BeNil)
	f.resolver = resolver
}

func (f *ResolverFixture) Teardown() {
	err := f.RemoveTempDir()
	f.So(err, should.BeNil)
}

func (f *ResolverFixture) TestResolveFile() {
	urlPath := "path/to/file.md"
	pathExpected, err := f.CreateFile(urlPath)
	f.So(err, should.BeNil)

	pathActual, err := f.resolver.Resolve(urlPath)
	f.So(err, should.BeNil)

	f.So(pathActual, should.Equal, pathExpected)
}

func (f *ResolverFixture) SkipTestResolveFileDecodesUrls() {
	// TODO: test spaces, high-unicode etc.
}

// ------------------------------------------------------------
// NewResolver

func TestNewResolver(t *testing.T) {
	gunit.Run(new(NewResolverFixture), t)
}

func (f *NewResolverFixture) Setup() {
	err := f.InitTempDir()
	f.So(err, should.BeNil)
}

func (f *NewResolverFixture) Teardown() {
	err := f.RemoveTempDir()
	f.So(err, should.BeNil)
}

type NewResolverFixture struct {
	*gunit.Fixture
	TempdirManager
}

func (f *NewResolverFixture) TestSetsAbsoluteRootDir() {
	rootDir, err := f.CreateDir("root")
	f.So(err, should.BeNil)
	rootDirAbs, err := filepath.Abs(rootDir)
	f.So(err, should.BeNil)

	r, err := NewResolver(rootDir)
	f.So(err, should.BeNil)
	f.So(r, should.NotBeNil)
	f.So(r.RootDir(), should.Equal, rootDirAbs)
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
	rootDir, err := f.CreateFile("root")
	f.So(err, should.BeNil)
	r, err := NewResolver(rootDir)
	f.So(r, should.BeNil)
	f.So(err, should.NotBeNil)
}

