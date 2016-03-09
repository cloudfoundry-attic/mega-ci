package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("main", func() {
	It("deploys manifest file specified", func() {
		var awsWasCalled bool
		awsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			awsWasCalled = true
		}))

		var boshcalls struct {
			DeleteCall int
			DeployCall int
			InfoCall   int
		}
		boshServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/deployments/multi-az-ssl":
				boshcalls.DeleteCall++
				w.Header().Set("Location", fmt.Sprintf("http://%s/tasks/1", r.Host))
				w.WriteHeader(http.StatusFound)
			case "/deployments":
				if r.Method == "POST" {
					boshcalls.DeployCall++
				}

				w.Header().Set("Location", fmt.Sprintf("http://%s/tasks/1", r.Host))
				w.WriteHeader(http.StatusFound)
			case "/tasks/1":
				w.Write([]byte(`{"state": "done"}`))
			case "/info":
				boshcalls.InfoCall++
				w.Write([]byte(`{"uuid":"some-director-uuid", "cpi":"some-cpi"}`))
			default:
				return
			}
		}))

		args := []string{
			"--manifest-path", "fixtures/multi-az-ssl.yml",
			"--director", boshServer.URL,
			"--user", "some-user",
			"--password", "some-password",
			"--aws-access-key-id", "some-aws-access-key-id",
			"--aws-secret-access-key", "some-aws-secret-access-key",
			"--aws-region", "some-aws-region",
			"--aws-endpoint-override", awsServer.URL,
		}

		session, err := gexec.Start(exec.Command(pathToMain, args...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		Expect(awsWasCalled).To(BeTrue())

		Expect(boshcalls.DeleteCall).To(Equal(1))
		Expect(boshcalls.DeployCall).To(Equal(1))
		Expect(boshcalls.InfoCall).To(Equal(1))
	})

	It("prints an error when the deployment fails", func() {
		awsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))

		boshServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))

		args := []string{
			"--manifest-path", "fixtures/multi-az-ssl.yml",
			"--director", boshServer.URL,
			"--user", "some-user",
			"--password", "some-password",
			"--aws-access-key-id", "some-aws-access-key-id",
			"--aws-secret-access-key", "some-aws-secret-access-key",
			"--aws-region", "some-aws-region",
			"--aws-endpoint-override", awsServer.URL,
		}

		session, err := gexec.Start(exec.Command(pathToMain, args...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))

		Expect(session.Err.Contents()).To(ContainSubstring("unexpected response 502 Bad Gateway"))
	})

	It("prints an error when an unknown flag is provided", func() {
		session, err := gexec.Start(exec.Command(pathToMain, "--some-unknown-flag"), GinkgoWriter, GinkgoWriter)

		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
		Expect(session.Err.Contents()).To(ContainSubstring("flag provided but not defined: -some-unknown-flag"))
	})
})
