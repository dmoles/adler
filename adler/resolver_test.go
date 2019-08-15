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
	tempdirFixture
	resolver *Resolver
}

func (f *ResolverFixture) SetupResolver() {
	resolver, err := NewResolver(f.tempDir)
	f.So(err, should.BeNil)
	f.resolver = resolver
}

func (f *ResolverFixture) TestResolveFile() {
	urlPath := "path/to/file.md"
	pathExpected := f.CreateFile(urlPath)
	pathActual, err := f.resolver.Resolve(urlPath)
	f.So(err, should.BeNil)
	f.So(pathActual, should.Equal, pathExpected)
}

func (f *ResolverFixture) SkipTestResolveFileDecodesUrls() {

}

// ------------------------------------------------------------
// NewResolver

func TestNewResolver(t *testing.T) {
	gunit.Run(new(NewResolverFixture), t)
}

type NewResolverFixture struct {
	tempdirFixture
}

func (f *NewResolverFixture) TestSetsAbsoluteRootDir() {
	rootDir := f.CreateDir("root")
	rootDirAbs, err := filepath.Abs(rootDir)
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
	rootDir := f.CreateFile("root")
	r, err := NewResolver(rootDir)
	f.So(r, should.BeNil)
	f.So(err, should.NotBeNil)
}

