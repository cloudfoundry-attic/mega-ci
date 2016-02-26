package awsdeployer

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
)

type SubnetChecker interface {
	CheckSubnets(manifestFilename string) (bool, error)
}

type AWSDeployer struct {
	bosh          clients.BOSH
	subnetChecker SubnetChecker
}

func NewAWSDeployer(bosh clients.BOSH, subnetChecker SubnetChecker) AWSDeployer {
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
	manifest, err := a.replaceUUID(manifestFilename)
	if err != nil {
		return err
	}

	buf, err := candiedyaml.Marshal(manifest)
	if err != nil {
		return err
	}

	err = a.bosh.Deploy(buf)
	if err != nil {
		return err
	}

	err = a.deleteDeployment(manifest)
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

func (a AWSDeployer) replaceUUID(manifestFilename string) (map[string]interface{}, error) {
	directorUUID, err := a.bosh.UUID()
	if err != nil {
		return nil, err
	}

	manifest, err := manifests.ReadManifest(manifestFilename)
	if err != nil {
		return nil, err
	}

	manifest["director_uuid"] = directorUUID

	return manifest, nil
}

func (a AWSDeployer) deleteDeployment(manifest map[string]interface{}) error {
	deploymentName, ok := manifest["name"].(string)
	if !ok {
		return errors.New("deployment name missing from manifest")
	}

	err := a.bosh.DeleteDeployment(deploymentName)
	if err != nil {
		return err
	}

	return nil
}
