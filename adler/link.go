package adler

import "os"

type Link interface {
	HREF() string
	Text() string
}

type Links []Link

func NewLink(fileInfo os.FileInfo) Link {
	return &link{fileInfo: fileInfo}
}

type link struct {
	fileInfo os.FileInfo
}

func (l *link) HREF() string {
	panic("implement me")
}

func (l *link) Text() string {
	panic("implement me")
}


