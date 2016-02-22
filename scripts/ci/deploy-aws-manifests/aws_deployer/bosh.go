package aws_deployer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type BOSH struct{}

func (b *BOSH) Deploy(manifest string, boshDirector string, boshUser string, boshPassword string) error {
	contents, err := ioutil.ReadFile(manifest)
	fmt.Println(string(contents))

	err = execute(os.Stdout, "-t", boshDirector, "-u", boshUser, "-p", boshPassword, "-d", manifest, "-n", "deploy")
	if err != nil {
		return err
	}

	return nil
}

func (b *BOSH) Status(boshDirector string, boshUser string, boshPassword string) (string, error) {
	output := new(bytes.Buffer)

	err := execute(output, "-t", boshDirector, "-u", boshUser, "-p", boshPassword, "status", "--uuid")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

func (b *BOSH) DeleteDeployment(deploymentName string, boshDirector string, boshUser string, boshPassword string) error {
	err := execute(os.Stdout, "-t", boshDirector, "-u", boshUser, "-p", boshPassword, "-n", "delete", "deployment", deploymentName)
	if err != nil {
		return err
	}

	return nil
}

func execute(output io.Writer, arguments ...string) error {
	boshBinary, err := exec.LookPath("bosh")
	if err != nil {
		return err
	}

	cmd := exec.Command(boshBinary, arguments...)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
