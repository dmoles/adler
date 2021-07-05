package markdown

// ------------------------------------------------------------
// Html implementation

func (b *mainContent) ToHtml() string {
	return b.html
}

// ------------------------------------------------------------
// Unexported

type mainContent struct {
	html string
}
