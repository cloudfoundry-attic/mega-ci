package deploy_cf_alone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCfAlone(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "deploy-cf-alone")
}
