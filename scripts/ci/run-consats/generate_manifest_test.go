package main

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("Generate", func() {
	var variables map[string]string

	BeforeEach(func() {
		variables = map[string]string{
			"AWS_AVAILIBILITY_ZONE":   "some-aws-availability-zone",
			"AWS_SUBNET_ID":           "some-aws-subnet-id",
			"AWS_ACCESS_KEY_ID":       "some-aws-access-key-id",
			"AWS_SECRET_ACCESS_KEY":   "some-aws-secret-access-key",
			"AWS_REGION":              "some-aws-region",
			"AWS_SECURITY_GROUP_NAME": "some-aws-security-group-name",
			"BOSH_DIRECTOR_UUID":      "some-bosh-director-uuid",
			"BOSH_TARGET":             "some-bosh-target",
			"BOSH_USERNAME":           "some-bosh-username",
			"BOSH_PASSWORD":           "some-bosh-password",
			"BOSH_DIRECTOR_CA_CERT":   "some-bosh-director-ca-cert",
			"REGISTRY_USERNAME":       "some-registry-username",
			"REGISTRY_PASSWORD":       "some-registry-password",
		}

		for name, value := range variables {
			variables[name] = os.Getenv(name)
			os.Setenv(name, value)
		}
	})

	AfterEach(func() {
		for name, value := range variables {
			os.Setenv(name, value)
		}
	})

	It("generates a manifest", func() {
		expectedManifest, err := ioutil.ReadFile("fixtures/expected.yml")
		Expect(err).NotTo(HaveOccurred())

		manifest, err := Generate("fixtures/example.yml")
		Expect(err).NotTo(HaveOccurred())

		Expect(manifest).To(MatchYAML(expectedManifest))
	})

	Context("failure cases", func() {
		It("returns an error when the example manifest does not exist", func() {
			_, err := Generate("fixtures/doesnotexist.yml")
			Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
		})

		It("returns an error when the example manifest is malformed", func() {
			_, err := Generate("fixtures/malformed.yml")
			Expect(err).To(MatchError(ContainSubstring("Invalid timestamp")))
		})
	})
})