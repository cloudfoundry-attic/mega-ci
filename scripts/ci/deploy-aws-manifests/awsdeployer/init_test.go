package awsdeployer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAWSDeployer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "awsdeployer")
}
