package markdown

import (
	"bufio"
	"bytes"
	"github.com/dmoles/adler/server/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
		extension.GFM,
		),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

const readmeMd = "README.md"

func DirToHTML(resolvedPath string, rootDir string) ([]byte, map[string]interface{}, error) {
	readmePath := filepath.Join(resolvedPath, readmeMd)
	if util.IsFile(readmePath) {
		return FileToHtml(readmePath)
	} else {
		return DirToIndexHtml(resolvedPath, rootDir)
	}
}

func DirToIndexHtml(dirPath string, rootDir string) ([]byte, Metadata, error) {
	dirIndex, err := NewDirIndex(dirPath)
	if err != nil {
		return nil, nil, err
	}
	dirIndexHtml, _, err := dirIndex.ToHtml(rootDir)
	if err != nil {
		return nil, nil, err
	}
	return dirIndexHtml, nil, nil
}

func GetTitle(in io.Reader) string {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()
		matches := headingRegexp.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func GetTitleFromFile(path string) (string, error) {
	if util.IsDirectory(path) {
		return AsTitle(path), nil
	}
	return ExtractTitle(path)
}

func AsTitle(path string) string {
	title := filepath.Base(path)
	title = strings.TrimSuffix(title, ".md")
	return strings.Title(title)
}

func ExtractTitle(path string) (string, error) {
	in, err := os.Open(path)
	defer util.CloseQuietly(in)
	if err != nil {
		return "", err
	}
	title := GetTitle(in)
	if title != "" {
		return title, nil
	}
	return AsTitle(path), nil
}

// ------------------------------------------------------------
// Unexported

var headingRegexp = regexp.MustCompile("^[\\s#]*#+ +(.+)$")

func stringToHtml(s string) ([]byte, Metadata, error) {
	return toHtml([]byte(s))
}

func toHtml(markdown []byte) ([]byte, Metadata, error) {
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(markdown, &buf, parser.WithContext(context)); err != nil {
		return nil, nil, err
	}
	return buf.Bytes(), meta.Get(context), nil
}
