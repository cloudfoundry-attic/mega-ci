package fakes

import (
	"github.com/pivotal-cf-experimental/bosh-test/bosh"
)

type BOSH struct {
	DeployCall struct {
		CallCount int
		Receives  struct {
			Manifest []byte
		}
		Returns struct {
			Error error
		}

		ReceivedManifests [][]byte
	}

	InfoCall struct {
		Returns struct {
			Info  bosh.DirectorInfo
			Error error
		}
	}

	DeleteDeploymentCall struct {
		Receives struct {
			Name string
		}

		Returns struct {
			Error error
		}
	}
}

func (b *BOSH) Deploy(manifest []byte) error {
	b.DeployCall.CallCount++
	b.DeployCall.Receives.Manifest = manifest
	b.DeployCall.ReceivedManifests = append(b.DeployCall.ReceivedManifests, manifest)
	return b.DeployCall.Returns.Error
}

func (b *BOSH) Info() (bosh.DirectorInfo, error) {
	return b.InfoCall.Returns.Info, b.InfoCall.Returns.Error
}

func (b *BOSH) DeleteDeployment(name string) error {
	b.DeleteDeploymentCall.Receives.Name = name
	return b.DeleteDeploymentCall.Returns.Error
}
