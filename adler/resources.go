package adler

import (
	"github.com/gobuffalo/packr"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TODO: use SCSS and integrate with a Mage build or something

var css = packr.NewBox("../css")

func findCSS(cssPath string) ([]byte, error) {
	return css.Find(cssPath)
}

var images = packr.NewBox("../images")

func findImage(imagePath string) ([]byte, string, error) {
	// TODO: straghten this out

	var data []byte
	var err error
	if filepath.IsAbs(imagePath) {
		data, err = ioutil.ReadFile(imagePath)
	} else {
		data, err = images.Find(imagePath)
	}

	if err != nil {
		return nil, "", err
	}
	if strings.HasSuffix(imagePath, ".png") {
		return data, "image/png", nil
	}
	if strings.HasSuffix(imagePath, ".ico") {
		return data, "image/vnd.microsoft.icon", nil
	}
	if strings.HasSuffix(imagePath, ".jpg") {
		return data, "image/jpeg", nil
	}
	return data, "application/octet-stream", nil
}