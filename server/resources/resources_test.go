package resources

import (
	"bytes"
	"crypto/md5"
	fmt "fmt"
	"github.com/dmoles/adler/server/util"
	"os"
	"path/filepath"
	"testing"
)

// ------------------------------------------------------------
// Helper functions

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
	resourcesDir := filepath.Join(util.ProjectRoot(), "resources")
	xBundle := newDirResources(resourcesDir)
	aBundle := defaultBundle

	_ = xBundle.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		if info.IsDir() {
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
	resourcesDir := filepath.Join(util.ProjectRoot(), "resources")
	expected := newDirResources(resourcesDir)
	actual := defaultBundle

	_ = actual.Walk(func(path string, info os.FileInfo, err error) error {
		_, err = expected.Get(path)
		if err != nil {
			t.Errorf("%v not found in %v, but present in %v: %v", path, expected, actual, err)
		}
		return nil
	})
}
