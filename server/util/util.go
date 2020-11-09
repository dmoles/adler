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

func ResolveRelative(urlPath string, rootDir string) (string, error) {
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
		log.Printf("can't stat %#v: %v", joinedPath, err)
		return "", errors.NotFound(joinedPath)
	}
	return joinedPath, nil
}

func ResolveDirectory(urlDirPath string, rootDir string) (string, error) {
	resolved, err := ResolveRelative(urlDirPath, rootDir)
	if err != nil {
		return "", err
	}
	return ToAbsoluteDirectory(resolved)
}

func ResolveFile(urlFilePath string, rootDir string) (string, error) {
	resolved, err := ResolveRelative(urlFilePath, rootDir)
	if err != nil {
		return "", err
	}
	return ToAbsoluteFile(resolved)
}

// TODO: use http.DetectContentType() instead?
func ContentType(urlPath string) string {
	if path.Base(urlPath) == "site.webmanifest" {
		return "application/manifest+json; charset=utf-8"
	}
	ext := path.Ext(urlPath)
	if ext == ".md" || ext == "" {
		return mime.TypeByExtension(".html")
	}
	if ext == ".ico" {
		return "image/x-icon"
	}
	if ext == ".woff" {
		return "font/woff"
	}
	if ext == ".woff2" {
		return "font/woff2"
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

func CloseQuietly(cl Closeable) {
	if cl != nil {
		err := cl.Close()
		if err != nil {
			msg := fmt.Sprintf("Error closing %v: %v\n", cl, err)
			log.Println(msg)
		}
	}
}

func RemoveAllQuietly(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		msg := fmt.Sprintf("Error removing %v: %v\n", path, err)
		log.Println(msg)
	}
}
