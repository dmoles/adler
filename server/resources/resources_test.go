package resources

import (
	"fmt"
	"github.com/dmoles/adler/server/util"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TODO: reverse test to make sure package doesn't include files deleted from the repo
func TestPackagedResources(t *testing.T) {
	expect := NewWithT(t).Expect

	resourcesDir := filepath.Join(util.ProjectRoot(), "resources")
	err := filepath.Walk(resourcesDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(resourcesDir, path)
		expect(err).NotTo(HaveOccurred())

		resourcePath := relativePath
		packagedFile, err := Open(resourcePath)
		expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error reading %s from resources: %v", resourcePath, err))

		pkgInfo, err := packagedFile.Stat()
		expect(err).NotTo(HaveOccurred())
		expect(pkgInfo.Name()).To(Equal(info.Name()))
		expect(pkgInfo.Size()).To(Equal(info.Size()))

		expectedData, err := ioutil.ReadFile(path)
		expect(err).NotTo(HaveOccurred())

		actualData, err := ioutil.ReadAll(packagedFile)
		expect(err).NotTo(HaveOccurred())

		expect(actualData).To(Equal(expectedData))

		//t.Logf("Verified %v = %v\n", resourcePath, path)

		return nil
	})
	expect(err).NotTo(HaveOccurred())
}
