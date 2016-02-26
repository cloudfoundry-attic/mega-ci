package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "main")
}

var (
	pathToMain string
)

var _ = BeforeSuite(func() {
	var err error

	pathToMain, err = gexec.Build("github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
