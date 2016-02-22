package manifests_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("manifests", func() {
	Describe("ReadManifest", func() {
		var (
			manifestsDirectory string
			err                error
		)

		BeforeEach(func() {
			manifestsDirectory, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the manifest as a map", func() {
			writeManifest(manifestsDirectory, "manifest.yml")

			manifestMap, err := manifests.ReadManifest(filepath.Join(manifestsDirectory, "manifest.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(manifestMap["director_uuid"]).To(Equal("BOSH-DIRECTOR-UUID"))
		})

		Context("failure cases", func() {
			It("returns an error when given invalid yaml", func() {
				manifestFile := filepath.Join(manifestsDirectory, "invalid_manifest.yml")
				err := ioutil.WriteFile(manifestFile, []byte("not: valid: yaml:"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifests.ReadManifest(manifestFile)
				Expect(err.Error()).To(ContainSubstring("mapping values are not allowed in this context"))
			})

			It("returns an error when the file doesn't exist", func() {
				_, err = manifests.ReadManifest("/nonexistent/file")
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})

	Describe("WriteManifest", func() {
		var (
			manifestsDirectory string
			err                error
		)

		BeforeEach(func() {
			manifestsDirectory, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
		})

		It("writes the given manifest to a file", func() {
			manifestToWrite := map[string]interface{}{"director_uuid": "BOSH-DIRECTOR-UUID"}
			manifestFile := filepath.Join(manifestsDirectory, "manifest.yml")

			err := manifests.WriteManifest(manifestFile, manifestToWrite)
			Expect(err).NotTo(HaveOccurred())

			writtenManifest := readManifest(manifestFile)
			Expect(writtenManifest["director_uuid"]).To(Equal("BOSH-DIRECTOR-UUID"))
		})

		Context("failure cases", func() {
			It("returns an error when manifest file cannot be written", func() {
				manifestToWrite := map[string]interface{}{"director_uuid": "BOSH-DIRECTOR-UUID"}

				err := manifests.WriteManifest("not/a/directory/manifest.yml", manifestToWrite)
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})
})

func writeManifest(directory string, filename string) {
	manifest := []byte("director_uuid: BOSH-DIRECTOR-UUID")
	err := ioutil.WriteFile(filepath.Join(directory, filename), manifest, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
}

func readManifest(manifestFile string) map[string]interface{} {
	file, err := os.Open(manifestFile)
	Expect(err).NotTo(HaveOccurred())
	defer file.Close()

	var document map[string]interface{}
	err = candiedyaml.NewDecoder(file).Decode(&document)
	Expect(err).NotTo(HaveOccurred())

	return document
}
