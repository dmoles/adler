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

func verify(expectedInfo os.FileInfo, expected Resources, actual Resources, path string) error {
	resourcePath := expected.RelativePath(path)

	actualInfo, err := Stat(actual, resourcePath)
	if err != nil {
		return fmt.Errorf("%v present in %v, but can't stat from %v: %v", resourcePath, expected, actual, err)
	}

	if expectedInfo.Name() != actualInfo.Name() {
		return fmt.Errorf("%v: expected %v, got %v", expectedInfo.Name(), expectedInfo.Name(), actualInfo.Name())
	}

	if expectedInfo.Size() != actualInfo.Size() {
		return fmt.Errorf("%v: expected size %v, got %v", expectedInfo.Name(), expectedInfo.Size(), actualInfo.Size())
	}

	expectedData, err := Read(expected, resourcePath)
	if err != nil {
		return fmt.Errorf("%v: can't read from %v: %v", resourcePath, expected, err)
	}
	actualData, err := Read(actual, resourcePath)
	if err != nil {
		return fmt.Errorf("%v: can't read from %v: %v", resourcePath, actual, err)
	}

	if !bytes.Equal(expectedData, actualData) {
		return fmt.Errorf("%v: expected %x (%d bytes), got %x (%d bytes)", resourcePath, md5.Sum(expectedData), len(expectedData), md5.Sum(actualData), len(actualData))
	}

	return nil
}

// ------------------------------------------------------------
// Tests

func TestPackagedResourcesMatchesResourceDir(t *testing.T) {
	resourcesDir := filepath.Join(util.ProjectRoot(), "resources")
	expected := newDirResources(resourcesDir)
	actual := defaultResources

	_ = expected.Walk(func(path string, expectedInfo os.FileInfo, err error) error {
		if err != nil {
			t.Error(err)
			return nil
		}
		if expectedInfo.IsDir() {
			return nil
		}
		err = verify(expectedInfo, expected, actual, path)
		if err != nil {
			t.Error(err)
		}
		return nil
	})
}

func TestPackagedResourcesHasNoExtraFiles(t *testing.T) {
	resourcesDir := filepath.Join(util.ProjectRoot(), "resources")
	expected := newDirResources(resourcesDir)
	actual := defaultResources

	_ = actual.Walk(func(resourcePath string, rightInfo os.FileInfo, err error) error {
		_, err = Stat(expected, resourcePath)
		if err != nil {
			t.Errorf("%v not found in %v, but present in %v: %v", resourcePath, expected, actual, err)
		}
		return nil
	})
}
