package manifests

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	Networks []Network
}

type Network struct {
	Subnets []struct {
		CloudProperties struct {
			Subnet string
		} `yaml:"cloud_properties"`
		Range string
	}
}

func ReadManifest(manifestFilename string) (map[string]interface{}, error) {
	contents, err := ioutil.ReadFile(manifestFilename)
	if err != nil {
		return nil, err
	}

	var document map[string]interface{}
	err = yaml.Unmarshal(contents, &document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func ReadNetworksFromManifest(manifestFilename string) ([]Network, error) {
	contents, err := ioutil.ReadFile(manifestFilename)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = yaml.Unmarshal(contents, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest.Networks, nil
}
