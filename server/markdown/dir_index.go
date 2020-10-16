package markdown

import (
	"fmt"
	"github.com/dmoles/adler/server/util"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ------------------------------------------------------------
// Exported

type DirIndex interface {
	ToHtml(rootDir string) ([]byte, error)
}

func NewDirIndex(dirPath string) (DirIndex, error) {
	pathsByTitle, err := getPathsByTitle(dirPath)
	if err != nil {
		return nil, err
	}
	return &dirIndex{
		dirPath:      dirPath,
		titles:       sortedTitles(pathsByTitle),
		pathsByTitle: pathsByTitle,
	}, nil
}

// ------------------------------------------------------------
// Unexported

type dirIndex struct {
	dirPath      string
	titles       []string
	pathsByTitle map[string]string
}

func (d *dirIndex) ToHtml(rootDir string) ([]byte, error) {
	title, err := GetTitleFromFile(d.dirPath)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	//noinspection GoUnhandledErrorResult
	fmt.Fprintf(&sb, "# %s\n\n", title)
	d.WriteMarkdown(&sb, rootDir)

	return stringToHtml(sb.String())
}

//noinspection GoUnhandledErrorResult
func (d *dirIndex) WriteMarkdown(w io.Writer, rootDir string) {
	for _, title := range d.titles {
		path := d.pathsByTitle[title]
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			log.Printf("Error determining relative path to file: %v: %v", path, err)
			continue
		}
		fmt.Fprintf(w, "- [%v](/%v)\n", title, relPath)
	}
}

// ------------------------------
// Utility methods

func sortedTitles(pathsByTitle map[string]string) []string {
	titles := make([]string, len(pathsByTitle))
	i := 0
	for k := range pathsByTitle {
		titles[i] = k
		i++
	}
	sort.Slice(titles, func(i, j int) bool {
		st1 := sortingTitle(titles[i])
		st2 := sortingTitle(titles[j])
		return st1 < st2
	})
	return titles
}

func getPathsByTitle(dirPath string) (map[string]string, error) {
	dirPath, err := util.ToAbsoluteDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	pathsByTitle := map[string]string{}
	for _, info := range files {
		filename := info.Name()
		if strings.HasPrefix(filename, ".") {
			continue
		}
		if !(info.IsDir() || strings.HasSuffix(filename, ".md")) {
			continue
		}
		fullPath := filepath.Join(dirPath, filename)
		title, err := GetTitleFromFile(fullPath)
		if err != nil {
			log.Printf("Error determining title from file: %v: %v", fullPath, err)
			continue
		}
		pathsByTitle[title] = fullPath
	}

	return pathsByTitle, nil
}

var numericPrefixRegexp = regexp.MustCompile("^[0-9-]+ (.+)")

func sortingTitle(t string) string {
	st := strings.TrimSpace(strings.ToLower(t))
	if submatch := numericPrefixRegexp.FindStringSubmatch(st); submatch != nil {
		st = submatch[1]
	}

	for _, prefix := range []string{"a ", "the "} {
		if strings.HasPrefix(st, prefix) {
			return strings.TrimPrefix(st, prefix)
		}
	}

	return st
}
