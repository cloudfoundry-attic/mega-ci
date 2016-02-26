package clients

import "github.com/pivotal-cf-experimental/bosh-test/bosh"

type BOSHClient interface {
	Deploy(manifest []byte) error
	DeleteDeployment(deploymentName string) error
	Info() (bosh.DirectorInfo, error)
}

type BOSH struct {
	boshClient BOSHClient
}

func NewBOSH(client BOSHClient) BOSH {
	return BOSH{
		boshClient: client,
	}
}

func (b BOSH) Deploy(manifest []byte) error {
	if err := b.boshClient.Deploy(manifest); err != nil {
		return err
	}
	return nil
}

func (b *BOSH) UUID() (string, error) {
	info, err := b.boshClient.Info()
	if err != nil {
		return "", err
	}

	return info.UUID, nil
}

func (b *BOSH) DeleteDeployment(deploymentName string) error {
	if err := b.boshClient.DeleteDeployment(deploymentName); err != nil {
		return err
	}

	return nil
}
