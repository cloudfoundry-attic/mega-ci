package subnetchecker_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/fakes"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/subnetchecker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubnetChecker", func() {
	var (
		manifestsDirectory string
		fakeAWS            *fakes.AWS
		subnetChecker      subnetchecker.SubnetChecker
	)

	BeforeEach(func() {
		fakeAWS = new(fakes.AWS)

		var err error
		manifestsDirectory, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
		subnetChecker = subnetchecker.NewSubnetChecker(fakeAWS)
	})

	It("returns true if all subnets in manifest exist on AWS", func() {
		const manifestWithSubnets1And2 = `---
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
		var awsSubnetsContainingSubnets1And2 = []clients.Subnet{
			{
				SubnetID:  "subnet-1",
				CIDRBlock: "10.0.20.0/24",
			},
			{
				SubnetID:  "subnet-2",
				CIDRBlock: "10.1.20.0/24",
			},
		}

		writeManifestWithBody(manifestsDirectory, "manifest.yml", manifestWithSubnets1And2)
		fakeAWS.FetchSubnetsCall.Returns.Subnets = awsSubnetsContainingSubnets1And2

		hasSubnets, err := subnetChecker.CheckSubnets(filepath.Join(manifestsDirectory, "manifest.yml"))

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeAWS.FetchSubnetsCall.Receives.SubnetIds).To(ConsistOf([]string{"subnet-1", "subnet-2"}))
		Expect(hasSubnets).To(BeTrue())
	})

	It("returns false if some subnet in the manifest does not exist on AWS", func() {
		const manifestWithSubnets1And2 = `---
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
		var awsSubnetsMissingSubnet2 = []clients.Subnet{
			{
				SubnetID:  "subnet-1",
				CIDRBlock: "10.0.20.0/24",
			},
		}

		writeManifestWithBody(manifestsDirectory, "manifest.yml", manifestWithSubnets1And2)
		fakeAWS.FetchSubnetsCall.Returns.Subnets = awsSubnetsMissingSubnet2

		hasSubnets, err := subnetChecker.CheckSubnets(filepath.Join(manifestsDirectory, "manifest.yml"))

		Expect(err).NotTo(HaveOccurred())
		Expect(hasSubnets).To(BeFalse())
		Expect(fakeAWS.FetchSubnetsCall.Receives.SubnetIds).To(ConsistOf([]string{"subnet-1", "subnet-2"}))
	})

	It("returns false if some subnet range in the manifest does not match the subnet range in AWS for the same id", func() {
		const manifestWithSubnets1And2 = `---
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
		var awsSubnetsMissingSubnet2 = []clients.Subnet{
			{
				SubnetID:  "subnet-1",
				CIDRBlock: "10.0.20.0/24",
			},
			{
				SubnetID:  "subnet-2",
				CIDRBlock: "10.3.20.0/24",
			},
		}

		writeManifestWithBody(manifestsDirectory, "manifest.yml", manifestWithSubnets1And2)
		fakeAWS.FetchSubnetsCall.Returns.Subnets = awsSubnetsMissingSubnet2

		hasSubnets, err := subnetChecker.CheckSubnets(filepath.Join(manifestsDirectory, "manifest.yml"))

		Expect(err).NotTo(HaveOccurred())
		Expect(hasSubnets).To(BeFalse())
		Expect(fakeAWS.FetchSubnetsCall.Receives.SubnetIds).To(ConsistOf([]string{"subnet-1", "subnet-2"}))
	})

	Context("failure cases", func() {

		It("returns an error when the manifest is not valid yaml", func() {
			writeManifestWithBody(manifestsDirectory, "invalid_manifest.yml", "not: valid: yaml:")
			_, err := subnetChecker.CheckSubnets(filepath.Join(manifestsDirectory, "invalid_manifest.yml"))
			Expect(err.Error()).To(ContainSubstring("mapping values are not allowed in this context"))
		})

		It("returns an error when aws client cannot get a session", func() {
			writeManifest(manifestsDirectory, "manifest.yml")
			fakeAWS.SessionCall.Returns.Error = errors.New("no aws session")

			_, err := subnetChecker.CheckSubnets(filepath.Join(manifestsDirectory, "manifest.yml"))
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("no aws session"))
		})

		It("returns an error when it FetchSubnets fails", func() {
			const manifestWithSubnets1And2 = `---
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
			writeManifestWithBody(manifestsDirectory, "manifest.yml", manifestWithSubnets1And2)
			fakeAWS.FetchSubnetsCall.Returns.Error = errors.New("something bad happened")

			_, err := subnetChecker.CheckSubnets(filepath.Join(manifestsDirectory, "manifest.yml"))
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("something bad happened"))
		})
	})
})

func writeManifest(directory string, filename string) {
	writeManifestWithBody(directory, filename, "director_uuid: BOSH-DIRECTOR-UUID\nname: a-deployment-name")
}

func writeManifestWithBody(directory string, filename string, body string) {
	err := ioutil.WriteFile(filepath.Join(directory, filename), []byte(body), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())
}
