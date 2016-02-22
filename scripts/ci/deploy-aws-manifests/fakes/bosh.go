package fakes

type Arguments struct {
	Manifest     string
	BoshDirector string
	BoshUser     string
	BoshPassword string
}

type BOSH struct {
	DeployCalls struct {
		Receives []Arguments
		Returns  struct {
			Error error
		}
	}

	StatusCall struct {
		Receives struct {
			BoshDirector string
			BoshUser     string
			BoshPassword string
		}

		Returns struct {
			UUID  string
			Error error
		}
	}

	DeleteDeploymentCall struct {
		Receives struct {
			DeploymentName string
			BoshDirector   string
			BoshUser       string
			BoshPassword   string
		}

		Returns struct {
			Error error
		}
	}
}

func (b *BOSH) Deploy(manifest string, boshDirector string, boshUser string, boshPassword string) error {
	b.DeployCalls.Receives = append(b.DeployCalls.Receives, Arguments{
		Manifest:     manifest,
		BoshDirector: boshDirector,
		BoshUser:     boshUser,
		BoshPassword: boshPassword,
	})

	return b.DeployCalls.Returns.Error
}

func (b *BOSH) Status(boshDirector string, boshUser string, boshPassword string) (string, error) {
	b.StatusCall.Receives.BoshDirector = boshDirector
	b.StatusCall.Receives.BoshUser = boshUser
	b.StatusCall.Receives.BoshPassword = boshPassword

	return b.StatusCall.Returns.UUID, b.StatusCall.Returns.Error
}

func (b *BOSH) DeleteDeployment(deploymentName string, boshDirector string, boshUser string, boshPassword string) error {
	b.DeleteDeploymentCall.Receives.DeploymentName = deploymentName
	b.DeleteDeploymentCall.Receives.BoshDirector = boshDirector
	b.DeleteDeploymentCall.Receives.BoshUser = boshUser
	b.DeleteDeploymentCall.Receives.BoshPassword = boshPassword

	return b.DeleteDeploymentCall.Returns.Error
}
