package deploy_etcd_cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func sourceCommand(command string, args ...string) *exec.Cmd {
	cmd := fmt.Sprintf(". task && %s %s", command, strings.Join(args, " "))
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
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	It("generates a config for the prepare-deployments tool", func() {
		deploymentsDir := tempDir + "/deployments-dir"

		err := os.Mkdir(tempDir+"/stemcell", os.ModePerm)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = os.Create(tempDir + "/stemcell/a-stemcell.tgz")
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
			CF:             tempDir + "/cf-release",
			ETCD:           tempDir + "/etcd-release",
			Stemcell:       tempDir + "/stemcell/a-stemcell.tgz",
			DeploymentsDir: deploymentsDir,
			Stubs:          []string{"/stub-1.yml", "/stub-2.yml"},
		}))
	})

	It("deploys the manifest with BOSH", func() {
		boshFile, err := os.Create(fmt.Sprintf("%s/bosh", tempDir))
		Expect(err).NotTo(HaveOccurred())

		err = boshFile.Chmod(os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		pathEnv := os.Getenv("PATH")
		os.Setenv("PATH", tempDir+":"+pathEnv)

		outputFile := tempDir + "/bosh-output"

		_, err = boshFile.WriteString(fmt.Sprintf("printf '%%s ' \"${@}\" >> %s", outputFile))
		Expect(err).NotTo(HaveOccurred())

		command := sourceCommand("deploy",
			"bosh.example.com",
			"username",
			"password",
			"manifest.yml")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		outputFileContents, err := ioutil.ReadFile(outputFile)
		Expect(err).NotTo(HaveOccurred())

		Expect(bytes.TrimSpace(outputFileContents)).To(Equal([]byte("-n -t bosh.example.com -u username -p password -d manifest.yml deploy")))
	})
})
