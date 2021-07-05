package markdown

import (
	"io/ioutil"
	"log"
)

// ------------------------------------------------------------
// Exported

type MarkdownFile interface {

}

func MarkdownFileFrom(resolvedPath string) (MarkdownFile, error) {
	md, bodyHtml, err := readMarkdown(resolvedPath)
	if err != nil {
		return nil, err
	}

	title := md.Title()
	if title == "" {
		// TODO: replace this with something that parses bodyHtml
		title, _ = ExtractTitle(resolvedPath)
	}

	mf := &markdownFile {
		title: title,
		bodyHtml:    bodyHtml,
		scripts:     md.Scripts(),
		stylesheets: md.Styles(),
	}
	return mf, nil
}



// ------------------------------------------------------------
// Unexported

type markdownFile struct {
	title string
	bodyHtml []byte
	stylesheets []Stylesheet
	scripts []Script
}

func (m *markdownFile) Stylesheets() []Stylesheet {
	return m.stylesheets
}

func (m *markdownFile) Scripts() []Script {
	return m.scripts
}

func readMarkdown(filePath string) (Metadata, []byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %v: %v", filePath, err)
		return nil, nil, err
	}

	htmlData, metadata, err := toHtml(data)
	if err != nil {
		log.Printf("Error parsing file %v: %v", filePath, err)
		return nil, nil, err
	}
	return metadata, htmlData, nil
}
