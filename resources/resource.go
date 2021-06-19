package resources

import (
	"fmt"
	"github.com/dmoles/adler/server/util"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// Resource An individual resource
type Resource interface {
	Path() string
	Bundle() Bundle
	Stat() os.FileInfo
	Open() (fs.File, error)
	Read() ([]byte, error)
	Copy(w io.Writer) (int64, error)
	Write(w http.ResponseWriter, urlPath string) error
	AsString() (string, error)
	Size() int64
	ContentType() string
}

func Resolve(resourceDir string, relativePath string) (Resource, error) {
	relativePathClean := path.Clean(relativePath)
	if relativePath == "" || strings.Contains(relativePath, "..") {
		return nil, fmt.Errorf("invalid resource path: %v", relativePath)
	}

	resolvedPath := path.Join(resourceDir, relativePathClean)
	return Get(resolvedPath)
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

func (r *resource) Open() (fs.File, error) {
	return r.bundle.Open(r.path)
}

func (r *resource) Read() ([]byte, error) {
	f, err := r.Open()
	if err != nil {
		return nil, err
	}
	defer util.CloseQuietly(f)
	return ioutil.ReadAll(f)
}

func (r *resource) Copy(w io.Writer) (int64, error) {
	f, err := r.Open()
	if err != nil {
		return 0, err
	}
	defer util.CloseQuietly(f)
	return io.Copy(w, f)
}

func (r *resource) Write(w http.ResponseWriter, urlPath string) error {
	size := r.Size()
	w.Header().Add("Content-Type", r.ContentType())
	w.Header().Add("Content-Length", strconv.FormatInt(size, 10))

	n, err := r.Copy(w)
	if n != size {
		if err == nil {
			return fmt.Errorf("wrote wrong number of bytes for %#v: expected %d, was %d", urlPath, size, n)
		}
		return fmt.Errorf("wrote wrong number of bytes for %#v: expected %d, was %d (%w)", urlPath, size, n, err)
	}
	return err
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
