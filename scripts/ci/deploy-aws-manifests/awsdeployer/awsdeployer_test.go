package awsdeployer_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/awsdeployer"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/fakes"
	"github.com/pivotal-cf-experimental/bosh-test/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWSDeployer", func() {
	Describe("Deploy", func() {
		var (
			manifestsDirectory string
			fakeBOSH           *fakes.BOSH
			fakeSubnetChecker  *fakes.SubnetChecker
			awsDeployer        awsdeployer.AWSDeployer
			stdout             io.Writer
		)

		BeforeEach(func() {
			fakeBOSH = &fakes.BOSH{}
			fakeSubnetChecker = &fakes.SubnetChecker{}
			stdout = ioutil.Discard

			var err error
			manifestsDirectory, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			awsDeployer = awsdeployer.NewAWSDeployer(clients.NewBOSH(fakeBOSH, stdout), fakeSubnetChecker, stdout)
			fakeSubnetChecker.CheckSubnetsCall.Returns.Bool = true
		})

		AfterEach(func() {
			os.RemoveAll(manifestsDirectory)
		})

		It("deploys all of the bosh manifests in the manifests directory", func() {
			writeManifestWithBody(manifestsDirectory, "first_manifest.yml", "director_uuid: BOSH-DIRECTOR-UUID\nname: deployment-1")
			writeManifestWithBody(manifestsDirectory, "second_manifest.yml", "director_uuid: BOSH-DIRECTOR-UUID\nname: deployment-2")
			writeManifest(manifestsDirectory, "not_a_manifest.json")

			deploymentError := awsDeployer.Deploy(manifestsDirectory)
			Expect(deploymentError).NotTo(HaveOccurred())

			Expect(fakeBOSH.DeployCall.CallCount).To(Equal(2))
			Expect(fakeBOSH.DeployCall.ReceivedManifests[0]).To(ContainSubstring("deployment-1"))
			Expect(fakeBOSH.DeployCall.ReceivedManifests[1]).To(ContainSubstring("deployment-2"))
		})

		It("deploys specified bosh manifest", func() {
			writeManifestWithBody(manifestsDirectory, "first_manifest.yml", "director_uuid: BOSH-DIRECTOR-UUID\nname: deployment-1")

			deploymentError := awsDeployer.Deploy(filepath.Join(manifestsDirectory, "first_manifest.yml"))
			Expect(deploymentError).NotTo(HaveOccurred())

			Expect(fakeBOSH.DeployCall.CallCount).To(Equal(1))
			Expect(fakeBOSH.DeployCall.ReceivedManifests[0]).To(ContainSubstring("deployment-1"))
		})

		It("replaces the bosh director uuid before deploying each manifest", func() {
			writeManifest(manifestsDirectory, "manifest.yml")
			fakeBOSH.InfoCall.Returns.Info = bosh.DirectorInfo{
				UUID: "retrieved-director-uuid",
			}
			deploymentError := awsDeployer.Deploy(manifestsDirectory)

			Expect(deploymentError).NotTo(HaveOccurred())

			Expect(string(fakeBOSH.DeployCall.ReceivedManifests[0])).To(ContainSubstring("director_uuid: retrieved-director-uuid"))
		})

		It("deletes the deployment", func() {
			writeManifestWithBody(manifestsDirectory, "manifest.yml", "director_uuid: BOSH-DIRECTOR-UUID\nname: some-deployment-name")

			deploymentError := awsDeployer.Deploy(manifestsDirectory)

			Expect(deploymentError).NotTo(HaveOccurred())
			Expect(fakeBOSH.DeleteDeploymentCall.Receives.Name).To(Equal("some-deployment-name"))
		})

		Context("failure cases", func() {
			It("returns an error when the BOSH manifest contains subnets not found on AWS", func() {
				const manifestWithSubnetC8 = `---
director_uuid: BOSH-DIRECTOR-UUID

name: multi-az-ssl

networks:
- subnets:
  - cloud_properties:
      subnet: "subnet-c8b76f90"
    range: 10.0.20.0/24
`
				writeManifestWithBody(manifestsDirectory, "manifest.yml", manifestWithSubnetC8)
				fakeSubnetChecker.CheckSubnetsCall.Returns.Bool = false

				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError).NotTo(BeNil())
				Expect(deploymentError.Error()).To(ContainSubstring("manifest subnets not found on AWS"))
			})

			It("returns an error when CheckSubnets returns an error", func() {
				writeManifest(manifestsDirectory, "manifest.yml")
				fakeSubnetChecker.CheckSubnetsCall.Returns.Error = errors.New("something bad happened")

				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError).NotTo(BeNil())
				Expect(deploymentError.Error()).To(ContainSubstring("something bad happened"))
			})

			It("returns an error when manifests directory does not exist", func() {
				deploymentError := awsDeployer.Deploy("/not/a/real/directory")
				Expect(deploymentError.Error()).To(ContainSubstring("no such file or directory"))
			})

			It("returns an error when bosh deploy fails", func() {
				writeManifest(manifestsDirectory, "manifest.yml")

				fakeBOSH.DeployCall.Returns.Error = errors.New("bosh deployment failed")

				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError.Error()).To(ContainSubstring("bosh deployment failed"))
			})

			It("returns an error when the manifest is not valid yaml", func() {
				writeManifestWithBody(manifestsDirectory, "invalid_manifest.yml", "not: valid: yaml:")
				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError.Error()).To(ContainSubstring("mapping values are not allowed in this context"))
			})

			It("returns an error when bosh UUID fails", func() {
				writeManifest(manifestsDirectory, "manifest.yml")

				fakeBOSH.InfoCall.Returns.Error = errors.New("bosh UUID failed")

				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError.Error()).To(ContainSubstring("bosh UUID failed"))
			})

			It("returns an error when the deployment name is not present in the manifest", func() {
				writeManifestWithBody(manifestsDirectory, "invalid_manifest.yml", "director_uuid: BOSH-DIRECTOR-UUID")
				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError.Error()).To(ContainSubstring("deployment name missing from manifest"))
			})

			It("returns an error when deletion of the deployment fails", func() {
				writeManifestWithBody(manifestsDirectory, "manifest.yml", "director_uuid: BOSH-DIRECTOR-UUID\nname: some-deployment-name")
				fakeBOSH.DeleteDeploymentCall.Returns.Error = errors.New("failed to delete deployment: some-deployment-name")

				deploymentError := awsDeployer.Deploy(manifestsDirectory)
				Expect(deploymentError.Error()).To(ContainSubstring("failed to delete deployment: some-deployment-name"))
			})
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
