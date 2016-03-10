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
			TaskId int
			Error  error
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

	GetTaskOutputCall struct {
		Receives struct {
			TaskId int
		}
		Returns struct {
			TaskOutputs []bosh.TaskOutput
			Error       error
		}
	}
}

func (b *BOSH) Deploy(manifest []byte) (int, error) {
	b.DeployCall.CallCount++
	b.DeployCall.Receives.Manifest = manifest
	b.DeployCall.ReceivedManifests = append(b.DeployCall.ReceivedManifests, manifest)
	return b.DeployCall.Returns.TaskId, b.DeployCall.Returns.Error
}

func (b *BOSH) Info() (bosh.DirectorInfo, error) {
	return b.InfoCall.Returns.Info, b.InfoCall.Returns.Error
}

func (b *BOSH) DeleteDeployment(name string) error {
	b.DeleteDeploymentCall.Receives.Name = name
	return b.DeleteDeploymentCall.Returns.Error
}

func (b *BOSH) GetTaskOutput(taskId int) ([]bosh.TaskOutput, error) {
	b.GetTaskOutputCall.Receives.TaskId = taskId
	return b.GetTaskOutputCall.Returns.TaskOutputs, b.GetTaskOutputCall.Returns.Error
}
