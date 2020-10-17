package resources

import (
	"fmt"
	"github.com/dmoles/adler/server/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// An individual resource
type Resource interface {
	Path() string
	Bundle() Bundle
	Stat() os.FileInfo
	Open() (http.File, error)
	Read() ([]byte, error)
	Copy(w io.Writer) (int64, error)
	AsString() (string, error)
	Size() int64
	ContentType() string
}

// ------------------------------------------------------------
// Unexported implementation

type resource struct {
	path   string
	bundle Bundle
	info   os.FileInfo
}

func (r *resource) Path() string {
	return r.path
}

func (r *resource) Bundle() Bundle {
	return r.bundle
}

func (r *resource) Stat() os.FileInfo {
	return r.info
}

func (r *resource) Open() (http.File, error) {
	return r.bundle.Open(r.path)
}

// TODO: consider caching data
func (r *resource) Read() ([]byte, error) {
	f, err := r.Open()
	if err != nil {
		return nil, err
	}
	defer util.CloseQuietly(f)
	return ioutil.ReadAll(f)
}

// TODO: consider caching data
func (r *resource) Copy(w io.Writer) (int64, error) {
	f, err := r.Open()
	if err != nil {
		return 0, err
	}
	defer util.CloseQuietly(f)
	return io.Copy(w, f)
}

func (r *resource) AsString() (string, error) {
	sb := new(strings.Builder)
	n, err := r.Copy(sb)
	if err != nil {
		return "", err
	}
	if n != r.Size() {
		return "", fmt.Errorf("%v: wrong number of bytes for Copy(); expected %v, was %v", r.path, r.Size(), n)
	}
	return sb.String(), nil
}

func (r *resource) Size() int64 {
	return r.info.Size()
}

func (r *resource) ContentType() string {
	return util.ContentType(r.path)
}
