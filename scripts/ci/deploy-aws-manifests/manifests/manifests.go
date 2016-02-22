package manifests

import (
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

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
