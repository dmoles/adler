package markdown

import (
	"bufio"
	"os"
	"regexp"

	"github.com/yuin/goldmark/util"

	. "github.com/dmoles/adler/server/util"
)

type Heading interface {
	Level() int
	Text() string
	Id() string
}

type heading struct {
	level int
	text  string
	id    string
}

func (h *heading) Level() int {
	return h.level
}

func (h *heading) Text() string {
	return h.text
}

func (h *heading) Id() string {
	return h.id
}

var headingRegexp = regexp.MustCompile("^\\s*(#+) +(.+)$")

func titleFromHeadings(headings []Heading) string {
	var current Heading
	for _, h := range headings {
		if current == nil || h.Level() < current.Level() {
			current = h
		}
	}
	if current == nil {
		return ""
	}
	return current.Text()
}

func findHeadings(filePath string) []Heading {
	in, err := os.Open(filePath)
	defer CloseQuietly(in)
	if err != nil {
		return nil
	}

	var headings []Heading

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()
		for _, m := range headingRegexp.FindAllStringSubmatch(text, -1) {
			h := headingFromSubmatch(m)
			if h != nil {
				headings = append(headings, h)
			}
		}
	}
	return headings
}

func headingFromSubmatch(m []string) *heading {
	if len(m) < 2 {
		return nil
	}

	lvl := len(m[1])
	text := m[2]

	return &heading{
		level: lvl,
		text:  text,
		id:    headingIdFrom(text),
	}
}

func headingIdFrom(text string) string {
	value := []byte(text)
	value = util.TrimLeftSpace(value)
	value = util.TrimRightSpace(value)
	result := []byte{}
	for i := 0; i < len(value); {
		v := value[i]
		l := util.UTF8Len(v)
		i += int(l)
		if l != 1 {
			continue
		}
		if util.IsAlphaNumeric(v) {
			if 'A' <= v && v <= 'Z' {
				v += 'a' - 'A'
			}
			result = append(result, v)
		} else if util.IsSpace(v) || v == '-' || v == '_' {
			result = append(result, '-')
		}
	}
	return string(value)
}
