package flags_test

import (
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/flags"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("flags", func() {
	It("extracts configuration data from the command line flags", func() {
		configuration, err := flags.ParseFlags([]string{
			"--manifests-directory", "some-manifests-directory",
			"--director", "some-director",
			"--user", "some-user",
			"--password", "some-password",
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(configuration.ManifestsDirectory).To(Equal("some-manifests-directory"))
		Expect(configuration.BoshDirector).To(Equal("some-director"))
		Expect(configuration.BoshUser).To(Equal("some-user"))
		Expect(configuration.BoshPassword).To(Equal("some-password"))
	})

	Describe("failure cases", func() {
		It("returns an error when flag parsing fails", func() {
			_, err := flags.ParseFlags([]string{"--wrong-flag", "some-string"})
			Expect(err.Error()).To(ContainSubstring("flag provided but not defined"))
		})
	})
})
