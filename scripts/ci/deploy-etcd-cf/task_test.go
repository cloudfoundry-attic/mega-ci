package deploy_etcd_cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func sourceCommand(command string, args ...string) *exec.Cmd {
	cmd := fmt.Sprintf("set -exu && . task && %s %s", command, strings.Join(args, " "))
	return exec.Command("bash", "-c", cmd)
}

type deploymentConfig struct {
	CF             string
	ETCD           string
	Stemcell       string
	DeploymentsDir string `json:"deployments-dir"`
	Stubs          []string
}

var _ = Describe("Task", func() {
	var (
		tempDir     string
		environment map[string]string
	)

	BeforeEach(func() {
		environment = map[string]string{
			"BOSH_DIRECTOR": os.Getenv("BOSH_DIRECTOR"),
			"BOSH_USER":     os.Getenv("BOSH_USER"),
			"BOSH_PASSWORD": os.Getenv("BOSH_PASSWORD"),
		}
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)

		for name, value := range environment {
			os.Setenv(name, value)
		}
	})

	It("generates a config for the prepare-deployments tool", func() {
		deploymentsDir := filepath.Join(tempDir, "/deployments-dir")

		err := os.Mkdir(filepath.Join(tempDir, "/stemcell"), os.ModePerm)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = os.Create(filepath.Join(tempDir, "/stemcell/a-stemcell.tgz"))
		Expect(err).ShouldNot(HaveOccurred())

		command := sourceCommand("generate_deployment_config",
			tempDir,
			deploymentsDir,
			"/stub-1.yml",
			"/stub-2.yml")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		config := deploymentConfig{}
		err = json.Unmarshal(session.Out.Contents(), &config)
		Expect(err).NotTo(HaveOccurred())

		Expect(config).To(Equal(deploymentConfig{
			CF:             "integration-latest",
			ETCD:           filepath.Join(tempDir, "/etcd-release"),
			Stemcell:       "integration-latest",
			DeploymentsDir: deploymentsDir,
			Stubs:          []string{"/stub-1.yml", "/stub-2.yml"},
		}))
	})

	It("deploys the manifest with BOSH", func() {
		boshFilePath := filepath.Join(tempDir, "bosh")

		pathEnv := os.Getenv("PATH")
		os.Setenv("PATH", tempDir+":"+pathEnv)

		outputFile := filepath.Join(tempDir, "/bosh-output")

		err := ioutil.WriteFile(
			boshFilePath,
			[]byte(fmt.Sprintf("printf '%%s ' \"${@}\" >> %s", outputFile)),
			os.ModePerm,
		)
		Expect(err).NotTo(HaveOccurred())

		command := sourceCommand("deploy", "bosh.example.com", "manifest.yml")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		outputFileContents, err := ioutil.ReadFile(outputFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(bytes.TrimSpace(outputFileContents)).
			To(Equal([]byte("-n -t bosh.example.com -d manifest.yml deploy --redact-diff")))
	})

	Describe("preflight_check", func() {
		Context("when the BOSH credentials are not set", func() {
			BeforeEach(func() {
				os.Setenv("BOSH_DIRECTOR", "")
				os.Setenv("BOSH_USER", "")
				os.Setenv("BOSH_PASSWORD", "")
			})

			It("fails the check", func() {
				command := sourceCommand("preflight_check")

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
			})
		})

		Context("when the BOSH credentials are set", func() {
			BeforeEach(func() {
				os.Setenv("BOSH_DIRECTOR", "bosh.example.com")
				os.Setenv("BOSH_USER", "username")
				os.Setenv("BOSH_PASSWORD", "password")
			})

			It("passes the check", func() {
				command := sourceCommand("preflight_check")

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(gexec.Exit(0))
			})
		})
	})
})
