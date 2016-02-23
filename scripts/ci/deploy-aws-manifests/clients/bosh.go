package clients

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type BOSH struct {
	director string
	user     string
	password string
}

func NewBOSH(boshDirector string, boshUser string, boshPassword string) *BOSH {
	return &BOSH{
		director: boshDirector,
		user:     boshUser,
		password: boshPassword,
	}
}

func (b *BOSH) Deploy(manifest string) error {
	contents, err := ioutil.ReadFile(manifest)
	fmt.Println(string(contents))

	err = execute(os.Stdout, "-t", b.director, "-u", b.user, "-p", b.password, "-d", manifest, "-n", "deploy")
	if err != nil {
		return err
	}

	return nil
}

func (b *BOSH) Status() (string, error) {
	output := new(bytes.Buffer)

	err := execute(output, "-t", b.director, "-u", b.user, "-p", b.password, "status", "--uuid")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.String()), nil
}

func (b *BOSH) DeleteDeployment(deploymentName string) error {
	err := execute(os.Stdout, "-t", b.director, "-u", b.user, "-p", b.password, "-n", "delete", "deployment", deploymentName)
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
