package main

import (
	"io/ioutil"
	"os"

	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {
	var variables map[string]string

	BeforeEach(func() {
		variables = map[string]string{
			"AWS_AVAILIBILITY_ZONE":   "some-aws-availability-zone",
			"AWS_SUBNETS":             `[{"id":"some-subnet-1","range":"10.0.4.0/24","az":"some-az-1"},{"id":"some-subnet-2","range":"10.0.5.0/24","az":"some-az-2"}]`,
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
			"PARALLEL_NODES":          "1",
			"CONSUL_RELEASE_VERSION":  "some-consul-release-version",
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

		Expect(manifest).To(gomegamatchers.MatchYAML(expectedManifest))
	})

	Context("failure cases", func() {
		It("returns an error when the parallel nodes is not an int", func() {
			os.Setenv("PARALLEL_NODES", "not an int")
			_, err := Generate("fixtures/example.yml")
			Expect(err).To(MatchError(ContainSubstring(`parsing "not an int": invalid syntax`)))
		})

		It("returns an error when the example manifest does not exist", func() {
			_, err := Generate("fixtures/doesnotexist.yml")
			Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
		})

		It("returns an error when the example manifest is malformed", func() {
			_, err := Generate("fixtures/malformed.yml")
			Expect(err).To(MatchError(ContainSubstring("cannot unmarshal !!str `this is...`")))
		})

		It("returns an error when the AWS_SUBNETS are not valid json", func() {
			os.Setenv("AWS_SUBNETS", "%%%%%")
			_, err := Generate("fixtures/example.yml")
			Expect(err).To(MatchError(ContainSubstring("invalid character '%' looking for beginning of value")))
		})
	})
})
