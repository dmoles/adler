package markdown

import "fmt"

type Stylesheet interface {
	ToHtml() string
}

type stylesheet struct {
	href string
}

func (s *stylesheet) ToHtml() string {
	return fmt.Sprintf("<link rel='stylesheet' href='%s'/>", s.href)
}
