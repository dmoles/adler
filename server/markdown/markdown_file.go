package markdown

// ------------------------------------------------------------
// Exported

type MarkdownFile interface {
	Title() Title
	Head() []Html
	MainContent() Html
	Headings() []Heading // TODO: use this to generate TOC
}

// ------------------------------------------------------------
// Unexported

type markdownFile struct {
	title       *title
	mainContent *mainContent
	stylesheets []*stylesheet
	scripts     []*script
	headings    []Heading
}

func fromParseResult(titleTxt string, mc *mainContent, md metadata, headings []Heading) *markdownFile {
	return &markdownFile{
		title:       &title{text: titleTxt},
		mainContent: mc,
		scripts:     md.Scripts(),
		stylesheets: md.Styles(),
		headings:    headings,
	}
}

func (m *markdownFile) Title() Title {
	return m.title
}

func (m *markdownFile) Head() []Html {
	var head []Html
	for _, ss := range m.stylesheets {
		head = append(head, ss)
	}
	for _, sc := range m.scripts {
		head = append(head, sc)
	}
	return head
}

func (m *markdownFile) MainContent() Html {
	return m.mainContent
}

func (m *markdownFile) Headings() []Heading {
	return m.headings
}
