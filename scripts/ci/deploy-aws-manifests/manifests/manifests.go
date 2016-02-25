package manifests

import (
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
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

func ReadManifest(manifestFile string) (map[string]interface{}, error) {
	file, err := os.Open(manifestFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var document map[string]interface{}
	err = candiedyaml.NewDecoder(file).Decode(&document)

	if err != nil {
		return nil, err
	}

	return document, nil
}

func WriteManifest(manifestFile string, document map[string]interface{}) error {
	fileToWrite, err := os.Create(manifestFile)
	defer fileToWrite.Close()
	if err != nil {
		return err
	}

	err = candiedyaml.NewEncoder(fileToWrite).Encode(document)
	if err != nil {
		return err
	}

	return nil
}

func ReadNetworksFromManifest(manifestFilename string) ([]Network, error) {
	file, err := os.Open(manifestFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var manifest Manifest
	err = candiedyaml.NewDecoder(file).Decode(&manifest)
	if err != nil {
		return nil, err
	}

	return manifest.Networks, nil
}
