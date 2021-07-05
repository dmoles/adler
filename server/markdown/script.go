package markdown

import "fmt"

type script struct {
	src string
	typ string
}

func (s *script) ToHtml() string {
	if s.typ == "" {
		return fmt.Sprintf("<script src=\"%s\"></script>", s.src)
	} else {
		return fmt.Sprintf("<script src=\"%s\" type=\"%s\"></script>", s.src, s.typ)
	}
}

func scriptFrom(md metadata) (*script, bool) {
	src := md.getString("src")
	if src != "" {
		s := script{
			src: src,
			typ: md.getString("type"),
		}
		return &s, true
	}
	return nil, false
}
