package clients_test

import (
	"bytes"
	"errors"
	"io/ioutil"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/fakes"
	"github.com/pivotal-cf-experimental/bosh-test/bosh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BOSH", func() {
	Describe("Deploy", func() {
		It("deploys a given manifest", func() {
			boshClient := &fakes.BOSH{}
			logger := bytes.NewBuffer([]byte{})
			client := clients.NewBOSH(boshClient, logger)
			boshClient.DeployCall.Returns.TaskId = 1
			boshClient.GetTaskOutputCall.Returns.TaskOutputs = []bosh.TaskOutput{
				{
					Time:     1457571106,
					Stage:    "some-stage",
					Task:     "some-task",
					State:    "some-state",
					Progress: 0,
				},
			}

			manifest := []byte("some-manifest")
			err := client.Deploy(manifest)
			Expect(err).NotTo(HaveOccurred())
			Expect(boshClient.DeployCall.ReceivedManifests[0]).To(Equal([]byte("some-manifest")))

			Expect(boshClient.GetTaskOutputCall.Receives.TaskId).To(Equal(1))
			Expect(logger.String()).To(ContainSubstring("Bosh Task 1:"))
			Expect(logger.String()).To(ContainSubstring("[Thu Mar 10 00:51:46 UTC 2016] Stage: some-stage Task: some-task State: some-state Progress: 0\n"))
		})

		Context("failure cases", func() {
			It("returns an error when the deployment fails", func() {
				boshClient := &fakes.BOSH{}
				client := clients.NewBOSH(boshClient, ioutil.Discard)
				boshClient.DeployCall.Returns.Error = errors.New("something bad happened")

				manifest := []byte("some-manifest")
				err := client.Deploy(manifest)
				Expect(err).To(MatchError("something bad happened"))
			})

			It("returns an error when get task output fails", func() {
				boshClient := &fakes.BOSH{}
				client := clients.NewBOSH(boshClient, ioutil.Discard)
				boshClient.GetTaskOutputCall.Returns.Error = errors.New("something bad happened")

				manifest := []byte("some-manifest")
				err := client.Deploy(manifest)
				Expect(err).To(MatchError("something bad happened"))
			})
		})
	})

	Describe("UUID", func() {
		It("retrieves the UUID from the bosh director", func() {
			boshClient := &fakes.BOSH{}
			boshClient.InfoCall.Returns.Info = bosh.DirectorInfo{
				UUID: "some-uuid",
			}
			client := clients.NewBOSH(boshClient, ioutil.Discard)

			uuid, err := client.UUID()
			Expect(err).NotTo(HaveOccurred())
			Expect(uuid).To(Equal("some-uuid"))
		})

		Context("failure cases", func() {
			It("returns an error when the deployment fails", func() {
				boshClient := &fakes.BOSH{}
				client := clients.NewBOSH(boshClient, ioutil.Discard)
				boshClient.InfoCall.Returns.Error = errors.New("something bad happened")

				_, err := client.UUID()
				Expect(err).To(MatchError("something bad happened"))
			})
		})
	})

	Describe("DeleteDeployment", func() {
		It("delete the deployment with the given name", func() {
			boshClient := &fakes.BOSH{}
			client := clients.NewBOSH(boshClient, ioutil.Discard)

			err := client.DeleteDeployment("some-deployment")
			Expect(err).NotTo(HaveOccurred())
			Expect(boshClient.DeleteDeploymentCall.Receives.Name).To(Equal("some-deployment"))
		})

		Context("failure cases", func() {
			It("returns an error when the deployment fails", func() {
				boshClient := &fakes.BOSH{}
				client := clients.NewBOSH(boshClient, ioutil.Discard)
				boshClient.DeleteDeploymentCall.Returns.Error = errors.New("something bad happened")

				err := client.DeleteDeployment("some-deployment-name")
				Expect(err).To(MatchError("something bad happened"))
			})
		})
	})
})
