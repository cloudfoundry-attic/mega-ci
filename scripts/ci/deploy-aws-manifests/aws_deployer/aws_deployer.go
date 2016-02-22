package aws_deployer

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
)

type AWSClient interface{}

type BOSHClient interface {
	Deploy(manifest string, boshDirector string, boshUser string, boshPassword string) error
	Status(boshDirector string, boshUser string, boshPassword string) (string, error)
}

type AWSDeployer struct {
	bosh BOSHClient
}

func NewAWSDeployer(aws AWSClient, bosh BOSHClient) AWSDeployer {
	return AWSDeployer{
		bosh: bosh,
	}
}

func (a AWSDeployer) Deploy(manifestsDirectory string, boshDirector string, boshUser string, boshPassword string) error {
	manifestsToDeploy, err := manifestsInDirectory(manifestsDirectory)
	if err != nil {
		return err
	}

	for _, manifest := range manifestsToDeploy {
		directorUUID, err := a.bosh.Status(boshDirector, boshUser, boshPassword)
		if err != nil {
			return err
		}

		err = replaceUUID(manifest, directorUUID)
		if err != nil {
			return err
		}

		err = a.bosh.Deploy(manifest, boshDirector, boshUser, boshPassword)
		if err != nil {
			return err
		}
	}

	return nil
}

func replaceUUID(manifestFile string, directorUUID string) error {
	document, err := manifests.ReadManifest(manifestFile)
	if err != nil {
		return err
	}

	document["director_uuid"] = directorUUID
	manifests.WriteManifest(manifestFile, document)

	return nil
}

func manifestsInDirectory(directory string) ([]string, error) {
	var manifests []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yml" {
			manifests = append(manifests, path)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return manifests, nil
}
