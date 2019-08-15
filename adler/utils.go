package adler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// Utility globals

var headingRegexp = regexp.MustCompile("#+ +(.+)$")

// ------------------------------------------------------------
// Utility methods

func closeQuietly(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing %v: %v\n", file.Name(), err)
	}
}

func invalidPath(urlPath string) error {
	return fmt.Errorf("invalid path: %#v", urlPath)
}

// Deprecated TODO: do we need this? does it work?
func textOfFirstHeading(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer closeQuietly(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := headingRegexp.FindStringSubmatch(scanner.Text())
		if len(matches) > 1 {
			return matches[1], nil
		}
	}
	return "", nil
}

func isMarkdownFile(name string) bool {
	return strings.HasSuffix(name, ".md")
}

// TODO: wrap os.Stat and ioutil.ReadDir to return something encapsulating full path
func relativeLink(parent string, info os.FileInfo) (string, bool) {
	name := info.Name()
	title := name
	relPath := name

	if !info.IsDir() {
		if !strings.HasSuffix(name, ".md") {
			return "", false
		}
		filePath := filepath.Join(parent, name)
		firstHeading, _ := textOfFirstHeading(filePath)
		if firstHeading != "" {
			title = firstHeading
		}
	}

	return fmt.Sprintf("[%v](%v)\n", title, relPath), true
}