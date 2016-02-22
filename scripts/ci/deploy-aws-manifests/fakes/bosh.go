package fakes

type BOSH struct {
	DeployCalls struct {
		Receives struct {
			Manifest []string
		}
		Returns struct {
			Error error
		}
	}

	StatusCall struct {
		Returns struct {
			UUID  string
			Error error
		}
	}

	DeleteDeploymentCall struct {
		Receives struct {
			DeploymentName string
		}

		Returns struct {
			Error error
		}
	}
}

func (b *BOSH) Deploy(manifest string) error {
	b.DeployCalls.Receives.Manifest = append(b.DeployCalls.Receives.Manifest, manifest)
	return b.DeployCalls.Returns.Error
}

func (b *BOSH) Status() (string, error) {
	return b.StatusCall.Returns.UUID, b.StatusCall.Returns.Error
}

func (b *BOSH) DeleteDeployment(deploymentName string) error {
	b.DeleteDeploymentCall.Receives.DeploymentName = deploymentName
	return b.DeleteDeploymentCall.Returns.Error
}
