package markdown

import "fmt"

// ------------------------------------------------------------
// Exported

type Title interface {
	ToHtml() string
	Text() string
}

// ------------------------------------------------------------
// Unexported

type title struct {
	text string
}

func (t *title) ToHtml() string {
	return fmt.Sprintf("<title>%s</title>", t.text)
}

func (t *title) Text() string {
	return t.text
}
