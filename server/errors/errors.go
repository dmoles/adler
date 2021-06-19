package errors

import "fmt"

func InvalidPath(urlPath string) error {
	return fmt.Errorf("invalid path: %#v", urlPath)
}
