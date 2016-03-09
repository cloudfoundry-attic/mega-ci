package flags

import (
	"flag"
	"io/ioutil"
)

type Configuration struct {
	ManifestPath        string
	BoshDirector        string
	BoshUser            string
	BoshPassword        string
	AWSAccessKeyID      string
	AWSSecretAccessKey  string
	AWSRegion           string
	AWSEndpointOverride string
}

func ParseFlags(arguments []string) (Configuration, error) {
	flags := flag.NewFlagSet("boshflags", flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	configuration := Configuration{}
	flags.StringVar(&configuration.ManifestPath, "manifest-path", "", "path to manifests directory or file")
	flags.StringVar(&configuration.BoshDirector, "director", "", "bosh director")
	flags.StringVar(&configuration.BoshUser, "user", "", "bosh user")
	flags.StringVar(&configuration.BoshPassword, "password", "", "bosh password")
	flags.StringVar(&configuration.AWSAccessKeyID, "aws-access-key-id", "", "aws access key id")
	flags.StringVar(&configuration.AWSSecretAccessKey, "aws-secret-access-key", "", "aws secret access key")
	flags.StringVar(&configuration.AWSRegion, "aws-region", "", "aws region")
	flags.StringVar(&configuration.AWSEndpointOverride, "aws-endpoint-override", "", "aws endpoint override")

	err := flags.Parse(arguments)
	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}
