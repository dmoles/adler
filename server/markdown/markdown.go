package markdown

import (
	"bufio"
	"bytes"
	"github.com/dmoles/adler/server/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// TODO: accept a Writer
func toHtml(markdown []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := md.Convert(markdown, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func FileToHtml(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %v: %v", filePath, err)
	}

	html, err := toHtml(data)
	if err != nil {
		log.Printf("Error parsing file %v: %v", filePath, err)
	}
	return html
}

func DirToHtml(dirPath string, rootDir string) ([]byte, error) {
	dirIndex, err := NewDirIndex(dirPath)
	if err != nil {
		return nil, err
	}
	dirIndexHtml, err := dirIndex.ToHtml(rootDir)
	if err != nil {
		return nil, err
	}
	return dirIndexHtml, nil
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
		return asTitle(path), nil
	}
	in, err := os.Open(path)
	defer util.CloseQuietly(in)
	if err != nil {
		return "", err
	}
	title := GetTitle(in)
	if title != "" {
		return title, nil
	}
	return asTitle(path), nil
}

func asTitle(filePath string) string {
	title := filepath.Base(filePath)
	title = strings.TrimSuffix(title, ".md")
	return strings.Title(title)
}

// ------------------------------------------------------------
// Unexported

var headingRegexp = regexp.MustCompile("^[\\s#]*#+ +(.+)$")

func stringToHtml(s string) ([]byte, error) {
	return toHtml([]byte(s))
}
