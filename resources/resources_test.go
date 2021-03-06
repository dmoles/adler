package resources

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/get-woke/go-gitignore"

	"github.com/dmoles/adler/server/util"
)

// ------------------------------------------------------------
// TODO: Something smarter

const resourceRoot = "resources"

// ------------------------------------------------------------
// Helper functions

var gitIgnore = func() *ignore.GitIgnore {
	path := filepath.Join(util.ProjectRoot(), ".gitignore")
	gi, err := ignore.CompileIgnoreFile(path)
	if err != nil {
		panic(err)
	}
	return gi
}()

func ignored(path string) bool {
	if filepath.IsAbs(path) {
		relativePath, err := filepath.Rel(util.ProjectRoot(), path)
		if err != nil {
			panic(err)
		}
		path = relativePath
	}
	// TODO: figure out how embeds really work so we can get rid of this
	if strings.HasSuffix(path, ".go") {
		return true
	}
	return gitIgnore.MatchesPath(path)
}

func verify(expected Resource, xBundle Bundle, aBundle Bundle, path string) error {
	resourcePath := xBundle.RelativePath(path)

	actual, err := aBundle.Get(resourcePath)
	if err != nil {
		return fmt.Errorf("%v present in %v, but can't get resource from %v: %v", resourcePath, xBundle, aBundle, err)
	}

	actualInfo := actual.Stat()
	expectedInfo := expected.Stat()
	if actualInfo.Name() != expectedInfo.Name() {
		return fmt.Errorf("%v: xBundle %v, got %v", expectedInfo.Name(), expectedInfo.Name(), actualInfo.Name())
	}
	if actualInfo.Size() != expectedInfo.Size() {
		return fmt.Errorf("%v: xBundle size %v, got %v", expectedInfo.Name(), expectedInfo.Size(), actualInfo.Size())
	}

	expectedData, err := expected.Read()
	if err != nil {
		return fmt.Errorf("%v: can't read from %v: %v", resourcePath, xBundle, err)
	}
	actualData, err := actual.Read()
	if err != nil {
		return fmt.Errorf("%v: can't read from %v: %v", resourcePath, aBundle, err)
	}

	if !bytes.Equal(expectedData, actualData) {
		return fmt.Errorf("%v: xBundle %x (%d bytes), got %x (%d bytes)", resourcePath, md5.Sum(expectedData), len(expectedData), md5.Sum(actualData), len(actualData))
	}

	return nil
}

// ------------------------------------------------------------
// Tests

func TestPackagedResourcesMatchesResourceDir(t *testing.T) {
	resourcesDir := filepath.Join(util.ProjectRoot(), resourceRoot)
	xBundle := newDirBundle(resourcesDir)
	aBundle := defaultBundle

	_ = xBundle.Walk(func(path string, d os.DirEntry, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		if d.IsDir() || ignored(path) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			t.Error(err)
			return nil
		}

		// TODO: something less awkward; also, validate paths
		expected := &resource{xBundle.RelativePath(path), xBundle, info}
		err = verify(expected, xBundle, aBundle, path)
		if err != nil {
			t.Error(err)
		}
		return nil
	})
}

func TestPackagedResourcesHasNoExtraFiles(t *testing.T) {
	resourcesDir := filepath.Join(util.ProjectRoot(), resourceRoot)
	expected := newDirBundle(resourcesDir)
	actual := defaultBundle

	_ = actual.Walk(func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		_, err = expected.Get(path)
		if err != nil {
			t.Errorf("%v not found in %v, but present in %v: %v", path, expected, actual, err)
		}
		return nil
	})
}
