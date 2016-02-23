package main

import (
	"os"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/aws_deployer"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/flags"
)

func main() {
	configuration, err := flags.ParseFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	aws := &aws_deployer.AWS{}
	bosh := aws_deployer.NewBOSH(configuration.BoshDirector, configuration.BoshUser, configuration.BoshPassword)

	awsDeployer := aws_deployer.NewAWSDeployer(
		aws,
		bosh,
	)

	err = awsDeployer.Deploy(configuration.ManifestsDirectory)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
