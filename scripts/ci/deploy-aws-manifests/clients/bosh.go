package clients

import (
	"fmt"
	"io"
	"time"

	"github.com/pivotal-cf-experimental/bosh-test/bosh"
)

type BOSHClient interface {
	Deploy(manifest []byte) (int, error)
	DeleteDeployment(deploymentName string) error
	GetTaskOutput(taskId int) ([]bosh.TaskOutput, error)
	Info() (bosh.DirectorInfo, error)
}

type BOSH struct {
	boshClient BOSHClient
	Logger     io.Writer
}

func NewBOSH(client BOSHClient, logger io.Writer) BOSH {
	return BOSH{
		boshClient: client,
		Logger:     logger,
	}
}

func (b BOSH) Deploy(manifest []byte) error {
	taskId, err := b.boshClient.Deploy(manifest)
	if err != nil {
		return err
	}

	taskOutputs, err := b.boshClient.GetTaskOutput(taskId)
	if err != nil {
		return err
	}

	fmt.Fprintf(b.Logger, "Bosh Task %d:\n", taskId)
	for _, taskOutput := range taskOutputs {
		fmt.Fprintf(b.Logger, "[%s] Stage: %s Task: %s State: %s Progress: %d\n", time.Unix(taskOutput.Time, 0).UTC().Format(time.UnixDate), taskOutput.Stage, taskOutput.Task, taskOutput.State, taskOutput.Progress)
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
