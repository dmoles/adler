package adler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type TempdirManager struct {
	tempDir string
}

func (t *TempdirManager) InitTempDir() error {
	tempDir, err := ioutil.TempDir("", "TempdirManager")
	if err != nil {
		return err
	}
	t.tempDir = tempDir
	return nil
}

func (t *TempdirManager) RemoveTempDir() error {
	return os.RemoveAll(t.tempDir)
}

// Create a subdirectory relative to the fixture's temporary
// directory, and return its path as a string
func (t *TempdirManager) CreateDir(relativeDir string) (string, error) {
	if t.tempDir == "" {
		return "", fmt.Errorf("TempdirManager: tempdir not initialized")
	}
	dirActual := path.Join(t.tempDir, relativeDir)
	err := os.MkdirAll(dirActual, 0755)
	if err != nil {
		return "", err
	}
	return dirActual, nil
}

// Create an empty file at a path relative to the fixture's temporary
// directory, and return its path as a string
func (t *TempdirManager) CreateFile(relativePath string) (string, error) {
	if t.tempDir == "" {
		return "", fmt.Errorf("TempdirManager: tempdir not initialized")
	}
	relativeDir := path.Dir(relativePath)
	dirActual, err := t.CreateDir(relativeDir)
	if err != nil {
		return "", err
	}

	baseName := path.Base(relativePath)
	pathActual := path.Join(dirActual, baseName)

	file, err := os.Create(pathActual)
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	return pathActual, nil
}