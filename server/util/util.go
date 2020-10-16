package util

import (
	"fmt"
	"github.com/dmoles/adler/server/errors"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func ToAbsoluteDirectory(dirPath string) (string, error) {
	dirPathAbs, err := filepath.Abs(dirPath)
	if err != nil {
		return "", nil
	}
	if IsDirectory(dirPathAbs) {
		return dirPathAbs, nil
	}
	return "", fmt.Errorf("not a directory: %s", dirPathAbs)
}

func IsDirectory(dirPath string) bool {
	f, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func IsFile(dirPath string) bool {
	f, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func ToAbsoluteFile(filePath string) (string, error) {
	filePathAbs, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}
	f, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	if f.IsDir() {
		return "", fmt.Errorf("not a plain file: %s", filePathAbs)
	}
	return filePathAbs, nil
}

// TODO: recreate a Resolver object and move these to it

func ResolvePath(urlPath string, rootDir string) (string, error) {
	decodedPath, err := url.PathUnescape(urlPath)
	if err != nil {
		return "", errors.InvalidPath(urlPath)
	}
	pathElements := strings.Split(decodedPath, "/")
	for _, pathElement := range pathElements {
		if pathElement == ".." {
			return "", errors.InvalidPath(urlPath)
		}
	}
	joinedPath := filepath.Join(rootDir, decodedPath)
	_, err = os.Stat(joinedPath)
	if err != nil {
		return "", err
	}
	return joinedPath, nil
}

func ResolveDirectory(urlDirPath string, rootDir string) (string, error) {
	resolved, err := ResolvePath(urlDirPath, rootDir)
	if err != nil {
		return "", err
	}
	return ToAbsoluteDirectory(resolved)
}

func ResolveFile(urlFilePath string, rootDir string) (string, error) {
	resolved, err := ResolvePath(urlFilePath, rootDir)
	if err != nil {
		return "", err
	}
	return ToAbsoluteFile(resolved)
}

// TODO: use http.DetectContentType() instead?
func ContentType(urlPath string) string {
	ext := path.Ext(urlPath)
	if ext == ".md" || ext == "" {
		return "text/html; charset=utf-8"
	}
	ct := mime.TypeByExtension(ext)
	if ct == "" {
		return "application/octet-stream"
	}
	return ct
}

func ContentLength(data []byte) string {
	return strconv.Itoa(len(data))
}

func WriteData(w http.ResponseWriter, urlPath string, data []byte) {
	w.Header().Add("Content-Type", ContentType(urlPath))
	w.Header().Add("Content-Length", ContentLength(data))
	n, err := w.Write(data)
	if err != nil {
		log.Printf("Error serving %#v: %v", urlPath, err)
	}
	if n != len(data) {
		log.Printf("Wrote wrong number of bytes for %#v: expected %d, was %d", urlPath, len(data), n)
	}
}

type Closeable interface {
	Close() error
}

func CloseQuietly(cl Closeable) func() {
	return func() {
		if cl != nil {
			err := cl.Close()
			if err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}
	}
}
