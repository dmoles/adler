package markdown

import (
	"bufio"
	"bytes"
	"github.com/dmoles/adler/server/util"
	"github.com/russross/blackfriday/v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FileToHtml(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %v: %v", filePath, err)
	}
	return blackfriday.Run(data)
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

func GetTitleFromString(s string) string {
	return GetTitle(bytes.NewBuffer([]byte(s)))
}

func GetTitleFromBytes(b []byte) string {
	return GetTitle(bytes.NewBuffer(b))
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
		return title, nil;
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
