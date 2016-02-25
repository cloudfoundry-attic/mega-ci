package awsdeployer

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
)

type BOSHClient interface {
	Deploy(manifest string) error
	Status() (string, error)
	DeleteDeployment(deploymentName string) error
}

type SubnetChecker interface {
	CheckSubnets(manifestFilename string) (bool, error)
}

type AWSDeployer struct {
	bosh          BOSHClient
	subnetChecker SubnetChecker
}

func NewAWSDeployer(bosh BOSHClient, subnetChecker SubnetChecker) AWSDeployer {
	return AWSDeployer{
		bosh:          bosh,
		subnetChecker: subnetChecker,
	}
}

func (a AWSDeployer) Deploy(manifestsDirectory string) error {
	manifestsToDeploy, err := manifestsInDirectory(manifestsDirectory)
	if err != nil {
		return err
	}

	for _, manifestFilename := range manifestsToDeploy {
		hasSubnets, err := a.subnetChecker.CheckSubnets(manifestFilename)
		if err != nil {
			return err
		}

		if !hasSubnets {
			return errors.New("manifest subnets not found on AWS")
		}

		err = a.deployManifest(manifestFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a AWSDeployer) deployManifest(manifestFilename string) error {
	err := a.replaceUUID(manifestFilename)
	if err != nil {
		return err
	}

	err = a.bosh.Deploy(manifestFilename)
	if err != nil {
		return err
	}

	err = a.deleteDeployment(manifestFilename)
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

func (a AWSDeployer) replaceUUID(manifestFilename string) error {
	directorUUID, err := a.bosh.Status()
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

func (a AWSDeployer) deleteDeployment(manifestFilename string) error {
	manifest, err := manifests.ReadManifest(manifestFilename)
	if err != nil {
		return err
	}

	deploymentName, ok := manifest["name"].(string)
	if !ok {
		return errors.New("deployment name missing from manifest")
	}

	err = a.bosh.DeleteDeployment(deploymentName)
	if err != nil {
		return err
	}

	return nil
}
