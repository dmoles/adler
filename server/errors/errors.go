package errors

import "fmt"

func InvalidPath(urlPath string) error {
	return fmt.Errorf("invalid path: %#v", urlPath)
}

func NotFound(urlPath string) error {
	return fmt.Errorf("not found: %#v", urlPath)
}
