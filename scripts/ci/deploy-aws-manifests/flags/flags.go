package flags

import (
	"flag"
	"io/ioutil"
)

type Configuration struct {
	ManifestsDirectory string
	BoshDirector       string
	BoshUser           string
	BoshPassword       string
}

func ParseFlags(arguments []string) (Configuration, error) {
	flags := flag.NewFlagSet("boshflags", flag.ContinueOnError)
	flags.SetOutput(ioutil.Discard)

	configuration := Configuration{}
	flags.StringVar(&configuration.ManifestsDirectory, "manifests-directory", "", "manifests directory")
	flags.StringVar(&configuration.BoshDirector, "director", "", "bosh director")
	flags.StringVar(&configuration.BoshUser, "user", "", "bosh user")
	flags.StringVar(&configuration.BoshPassword, "password", "", "bosh password")

	err := flags.Parse(arguments)
	if err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}
