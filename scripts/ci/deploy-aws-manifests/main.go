package main

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/aws_deployer"
)

func main() {
	aws := &aws_deployer.AWS{}
	bosh := aws_deployer.NewBOSH(os.Args[2], os.Args[3], os.Args[4])

	awsDeployer := aws_deployer.NewAWSDeployer(
		aws,
		bosh,
	)

	err := awsDeployer.Deploy(filepath.Join(os.Args[1], "manifests/aws"))
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
