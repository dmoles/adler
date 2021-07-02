package resources

import (
	"embed"
)

// ------------------------------------------------------------
// Exported

func Get(resourcePath string) (Resource, error) {
	return defaultBundle.Get(resourcePath)
}

// ------------------------------------------------------------
// Unexported

// TODO: figure out how resources work so we can do something a little less messy
//go:embed *.ico
//go:embed *.png
//go:embed *.webmanifest
//go:embed css
//go:embed fonts
//go:embed images
//go:embed templates
var defaultResources embed.FS

var defaultBundle = func() Bundle {
	return &bundle{"adler", defaultResources}
}()
