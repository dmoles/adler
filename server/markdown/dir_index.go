package markdown

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dmoles/adler/server/util"
)

// ------------------------------------------------------------
// Unexported

type dirIndex struct {
	title        string
	dirPath      string
	titles       []string
	pathsByTitle map[string]string
}

func newDirIndex(dirPath string) (*dirIndex, error) {
	var title string

	readmePath := filepath.Join(dirPath, readmeMd)
	if util.IsFile(readmePath) {
		mf, err := FromFile(readmePath)
		if err != nil {
			return nil, err
		}
		title = mf.Title().Text()
	}
	if title == "" {
		stem := filepath.Base(dirPath)
		title = strings.Title(stem)
	}

	pathsByTitle, err := getPathsByTitle(dirPath)
	if err != nil {
		return nil, err
	}

	return &dirIndex{
		title:        title,
		dirPath:      dirPath,
		titles:       sortedTitles(pathsByTitle),
		pathsByTitle: pathsByTitle,
	}, nil
}

func (d *dirIndex) toMarkdownFile(basePath string) (MarkdownFile, error) {
	dirPath := d.dirPath
	relPath := filepath.Base(dirPath)

	var sb strings.Builder
	//noinspection GoUnhandledErrorResult
	fmt.Fprintf(&sb, "# [%s](%s)\n\n", d.title, relPath)

	for _, title := range d.titles {
		path := d.pathsByTitle[title]
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			log.Printf("Error determining relative path to file: %v: %v", path, err)
			continue
		}
		//noinspection GoUnhandledErrorResult
		fmt.Fprintf(&sb, "- [%v](%v)\n", title, relPath)
	}

	mc, md, err := parseString(sb.String())
	if err != nil {
		return nil, err
	}
	return fromParseResult(d.title, mc, md, nil), nil
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
	// TODO: cache this so we're not constantly reparsing every file
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	pathsByTitle := map[string]string{}
	for _, info := range files {
		baseName := info.Name()
		if baseName == readmeMd || strings.HasPrefix(baseName, ".") {
			continue
		}
		if !(info.IsDir() || strings.HasSuffix(baseName, mdExt)) {
			continue
		}
		var mf MarkdownFile

		fullPath := filepath.Join(dirPath, baseName)
		if util.IsDirectory(fullPath) {
			mf, err = ForDirectory(fullPath)
		} else {
			mf, err = FromFile(fullPath)
		}
		if err != nil {
			log.Printf("Error determining title from file: %v: %v", fullPath, err)
			continue
		}
		title := mf.Title().Text()
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
