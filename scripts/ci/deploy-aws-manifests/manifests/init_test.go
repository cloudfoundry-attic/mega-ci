package manifests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestManifests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "manifests")
}
