package markdown

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/dmoles/adler/server/util"
)

// ------------------------------------------------------------
// Exported

// TODO: cache this so we don't keep parsing everything repeatedly
func FromFile(filePath string) (MarkdownFile, error) {
	mc, md, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	headings := findHeadings(filePath)

	titleTxt := md.Title()
	if titleTxt == "" {
		titleTxt = titleFromHeadings(headings)
	}
	if titleTxt == "" {
		baseName := filepath.Base(filePath)
		stem := strings.TrimSuffix(baseName, mdExt)
		if strings.ToUpper(stem) == readme {
			dirPath := filepath.Dir(filePath)
			stem = filepath.Base(dirPath)
		}
		titleTxt = strings.Title(stem)
	}

	return fromParseResult(titleTxt, mc, md, headings), nil
}

func ForDirectory(dirPath string) (MarkdownFile, error) {
	readmePath := filepath.Join(dirPath, readmeMd)
	if util.IsFile(readmePath) {
		return FromFile(readmePath)
	}

	return DirectoryIndex(dirPath, dirPath)
}

func DirectoryIndex(dirPath string, basePath string) (MarkdownFile, error) {
	dx, err := newDirIndex(dirPath)
	if err != nil {
		return nil, err
	}
	return dx.toMarkdownFile(basePath)
}

var md = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
		extension.GFM,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

// ------------------------------------------------------------
// Unexported

const readme = "README"
const mdExt = ".md"
const readmeMd = readme + mdExt

func parseFile(filePath string) (*mainContent, metadata, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %v: %v", filePath, err)
		return nil, nil, err
	}
	return parseBytes(data)
}

func parseString(markdown string) (*mainContent, metadata, error) {
	return parseBytes([]byte(markdown))
}

func parseBytes(markdown []byte) (*mainContent, metadata, error) {
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(markdown, &buf, parser.WithContext(context)); err != nil {
		return nil, nil, err
	}
	mainContentStr := string(buf.Bytes())
	return &mainContent{mainContentStr}, meta.Get(context), nil
}
