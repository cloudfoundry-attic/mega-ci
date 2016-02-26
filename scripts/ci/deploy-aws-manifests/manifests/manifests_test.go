package manifests_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const manifestWithSubnets = `---
director_uuid: BOSH-DIRECTOR-UUID

name: multi-az-ssl

networks:
- subnets:
  - cloud_properties:
      subnet: "subnet-1"
    range: 10.0.20.0/24
- subnets:
  - cloud_properties:
      subnet: "subnet-2"
    range: 10.1.20.0/24
`

var _ = Describe("manifests", func() {
	var (
		manifestsDirectory string
		err                error
	)

	BeforeEach(func() {
		manifestsDirectory, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ReadManifest", func() {
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

	Describe("ReadNetworksFromManifest", func() {
		It("reads the network information from the given manifest", func() {
			writeManifestWithBody(manifestsDirectory, "manifest-with-subnets.yml", manifestWithSubnets)

			networks, err := manifests.ReadNetworksFromManifest(filepath.Join(manifestsDirectory, "manifest-with-subnets.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(networks[0].Subnets[0].CloudProperties.Subnet).To(Equal("subnet-1"))
			Expect(networks[0].Subnets[0].Range).To(Equal("10.0.20.0/24"))

			Expect(networks[1].Subnets[0].CloudProperties.Subnet).To(Equal("subnet-2"))
			Expect(networks[1].Subnets[0].Range).To(Equal("10.1.20.0/24"))
		})
	})

})

func writeManifest(directory string, filename string) {
	manifest := []byte("director_uuid: BOSH-DIRECTOR-UUID")
	err := ioutil.WriteFile(filepath.Join(directory, filename), manifest, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
}

func writeManifestWithBody(directory string, filename string, body string) {
	err := ioutil.WriteFile(filepath.Join(directory, filename), []byte(body), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
}
