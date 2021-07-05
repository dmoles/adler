package markdown

import "fmt"

type stylesheet struct {
	href string
}

func (s *stylesheet) ToHtml() string {
	return fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"/>", s.href)
}
