package main

import (
	"os"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/awsdeployer"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/flags"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/subnetchecker"
)

func main() {
	configuration, err := flags.ParseFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	aws := clients.NewAWS(configuration.AWSAccessKeyID, configuration.AWSSecretAccessKey, configuration.AWSRegion)
	bosh := clients.NewBOSH(configuration.BoshDirector, configuration.BoshUser, configuration.BoshPassword)
	subnetChecker := subnetchecker.NewSubnetChecker(aws)

	awsDeployer := awsdeployer.NewAWSDeployer(bosh, subnetChecker)

	err = awsDeployer.Deploy(configuration.ManifestsDirectory)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
