package deploy_etcd_cf

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDeployEtcdCf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DeployEtcdCf Suite")
}
