package deploy_consul_cf

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDeployConsulCf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "deploy-consul-cf")
}
