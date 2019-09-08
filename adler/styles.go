package adler

import (
	"github.com/gobuffalo/packr"
)

// TODO: use SCSS and integrate with a Mage build or something

var css = packr.NewBox("../css")

func findCSS(cssPath string) ([]byte, error) {
	return css.Find(cssPath)
}