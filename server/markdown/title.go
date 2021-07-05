package markdown

import "fmt"

// ------------------------------------------------------------
// Exported

type Title interface {
	ToHtml() string
}

// ------------------------------------------------------------
// Unexported

type title struct {
	title string
}

func (t *title) ToHtml() string {
	return fmt.Sprintf()
}
