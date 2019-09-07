package adler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/lithammer/dedent"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ------------------------------------------------------------
// Utility globals

var headingRegexp = regexp.MustCompile("^[\\s#]*#+ +(.+)$")

// ------------------------------------------------------------
// Utility methods

func invalidPath(urlPath string) error {
	return fmt.Errorf("invalid path: %#v", urlPath)
}

func isMarkdownFile(name string) bool {
	return strings.HasSuffix(name, ".md")
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func isDirectory(dirPath string) bool {
	f, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func textOfFirstHeading(markdownBody []byte) string {
	scanner := bufio.NewScanner(bytes.NewBuffer(markdownBody))
	for scanner.Scan() {
		text := scanner.Text()
		matches := headingRegexp.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
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

func asTitle(filePath string) string {
	title := filepath.Base(filePath)
	if isMarkdownFile(title) {
		title = strings.TrimSuffix(title, ".md")
	}
	return strings.Title(title)
}