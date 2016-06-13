package awsdeployer_test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

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
			fakeBOSH          *fakes.BOSH
			fakeSubnetChecker *fakes.SubnetChecker
			awsDeployer       awsdeployer.AWSDeployer
			stdout            io.Writer
		)

		BeforeEach(func() {
			fakeBOSH = &fakes.BOSH{}
			fakeSubnetChecker = &fakes.SubnetChecker{}
			stdout = ioutil.Discard

			awsDeployer = awsdeployer.NewAWSDeployer(clients.NewBOSH(fakeBOSH, stdout), fakeSubnetChecker, stdout)
			fakeSubnetChecker.CheckSubnetsCall.Returns.Bool = true
		})

		It("deploys specified bosh manifest", func() {
			manifestFilename := createManifestFile("director_uuid: BOSH-DIRECTOR-UUID\nname: deployment-1")

			deploymentError := awsDeployer.Deploy(manifestFilename)
			Expect(deploymentError).NotTo(HaveOccurred())

			Expect(fakeBOSH.DeployCall.CallCount).To(Equal(1))
			Expect(fakeBOSH.DeployCall.Receives.Manifest).To(ContainSubstring("deployment-1"))
		})

		It("replaces the bosh director uuid before deploying each manifest", func() {
			manifestFilename := createBasicManifestFile()
			fakeBOSH.InfoCall.Returns.Info = bosh.DirectorInfo{
				UUID: "retrieved-director-uuid",
			}
			deploymentError := awsDeployer.Deploy(manifestFilename)

			Expect(deploymentError).NotTo(HaveOccurred())

			Expect(string(fakeBOSH.DeployCall.Receives.Manifest)).To(ContainSubstring("director_uuid: retrieved-director-uuid"))
		})

		It("deletes the deployment", func() {
			manifestFilename := createManifestFile("director_uuid: BOSH-DIRECTOR-UUID\nname: some-deployment-name")

			deploymentError := awsDeployer.Deploy(manifestFilename)

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

				manifestFilename := createManifestFile(manifestWithSubnetC8)
				fakeSubnetChecker.CheckSubnetsCall.Returns.Bool = false

				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("manifest subnets not found on AWS"))
			})

			It("returns an error when CheckSubnets returns an error", func() {
				manifestFilename := createBasicManifestFile()
				fakeSubnetChecker.CheckSubnetsCall.Returns.Error = errors.New("something bad happened")

				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("something bad happened"))
			})

			It("returns an error when bosh deploy fails", func() {
				manifestFilename := createBasicManifestFile()

				fakeBOSH.DeployCall.Returns.Error = errors.New("bosh deployment failed")

				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("bosh deployment failed"))
			})

			It("returns an error when the manifest is not valid yaml", func() {
				manifestFilename := createManifestFile("not: valid: yaml:")
				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("mapping values are not allowed in this context"))
			})

			It("returns an error when bosh UUID fails", func() {
				manifestFilename := createBasicManifestFile()

				fakeBOSH.InfoCall.Returns.Error = errors.New("bosh UUID failed")

				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("bosh UUID failed"))
			})

			It("returns an error when the deployment name is not present in the manifest", func() {
				manifestFilename := createManifestFile("director_uuid: BOSH-DIRECTOR-UUID")
				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("deployment name missing from manifest"))
			})

			It("returns an error when deletion of the deployment fails", func() {
				manifestFilename := createManifestFile("director_uuid: BOSH-DIRECTOR-UUID\nname: some-deployment-name")
				fakeBOSH.DeleteDeploymentCall.Returns.Error = errors.New("failed to delete deployment: some-deployment-name")

				deploymentError := awsDeployer.Deploy(manifestFilename)
				Expect(deploymentError.Error()).To(ContainSubstring("failed to delete deployment: some-deployment-name"))
			})
		})
	})
})

func createBasicManifestFile() string {
	return createManifestFile("director_uuid: BOSH-DIRECTOR-UUID\nname: a-deployment-name")
}

func createManifestFile(body string) string {
	manifest, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(manifest.Name(), []byte(body), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	return manifest.Name()
}
