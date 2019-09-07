package adler

import (
	"bufio"
	"fmt"
	"github.com/lithammer/dedent"
	"log"
	"os"
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

// Deprecated TODO: consider reading markdownPage.Content at initialization
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

// Dedents and trims whitespace from the specified string. Preserves
// up to 1 trailing newline.
func trim(text string) string {
	dedented := dedent.Dedent(text)
	trimmed := strings.TrimSpace(dedented)
	if strings.HasSuffix(dedented, "\n") {
		return trimmed + "\n"
	}
	return trimmed
}