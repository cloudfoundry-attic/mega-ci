package aws_deployer

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
)

type AWSClient interface{}

type BOSHClient interface {
	Deploy(manifest string, boshDirector string, boshUser string, boshPassword string) error
	Status(boshDirector string, boshUser string, boshPassword string) (string, error)
	DeleteDeployment(deploymentName string, boshDirector string, boshUser string, boshPassword string) error
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
		err := a.deployManifest(manifest, boshDirector, boshUser, boshPassword)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a AWSDeployer) deployManifest(manifestFilename string, boshDirector string, boshUser string, boshPassword string) error {
	err := a.replaceUUID(manifestFilename, boshDirector, boshUser, boshPassword)
	if err != nil {
		return err
	}

	err = a.bosh.Deploy(manifestFilename, boshDirector, boshUser, boshPassword)
	if err != nil {
		return err
	}

	err = a.deleteDeployment(manifestFilename, boshDirector, boshUser, boshPassword)
	if err != nil {
		return err
	}

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

func (a AWSDeployer) replaceUUID(manifestFilename string, boshDirector string, boshUser string, boshPassword string) error {
	directorUUID, err := a.bosh.Status(boshDirector, boshUser, boshPassword)
	if err != nil {
		return err
	}

	manifest, err := manifests.ReadManifest(manifestFilename)
	if err != nil {
		return err
	}

	manifest["director_uuid"] = directorUUID
	manifests.WriteManifest(manifestFilename, manifest)

	return nil
}

func (a AWSDeployer) deleteDeployment(manifestFilename string, boshDirector string, boshUser string, boshPassword string) error {
	manifest, err := manifests.ReadManifest(manifestFilename)
	if err != nil {
		return err
	}

	deploymentName, ok := manifest["name"].(string)
	if !ok {
		return errors.New("deployment name missing from manifest")
	}

	err = a.bosh.DeleteDeployment(deploymentName, boshDirector, boshUser, boshPassword)
	if err != nil {
		return err
	}

	return nil
}
