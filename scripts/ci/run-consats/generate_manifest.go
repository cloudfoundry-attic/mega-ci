package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Network struct {
	Name      string   `yaml:"name"`
	StaticIPs []string `yaml:"static_ips"`
}

type Manifest struct {
	Name         interface{} `yaml:"name"`
	DirectorUUID string      `yaml:"director_uuid"`
	Releases     []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"releases"`
	Stemcells []struct {
		Alias   string `yaml:"alias"`
		OS      string `yaml:"os"`
		Version string `yaml:"version"`
	} `yaml:"stemcells"`
	InstanceGroups []struct {
		Instances    int       `yaml:"instances"`
		Name         string    `yaml:"name"`
		Lifecycle    string    `yaml:"lifecycle"`
		VMExtensions []string  `yaml:"vm_extensions"`
		VMType       string    `yaml:"vm_type"`
		Stemcell     string    `yaml:"stemcell"`
		AZs          []string  `yaml:"azs"`
		Networks     []Network `yaml:"networks"`
		Jobs         []struct {
			Name    string `yaml:"name'`
			Release string `yaml:"release'`
		} `yaml:"jobs"`
	} `yaml:"instance_groups"`
	Properties struct {
		Consul struct {
			AcceptanceTests struct {
				AWS struct {
					AccessKeyID           string      `yaml:"access_key_id"`
					SecretAccessKey       string      `yaml:"secret_access_key"`
					Region                string      `yaml:"region"`
					DefaultKeyName        interface{} `yaml:"default_key_name"`
					DefaultSecurityGroups []string    `yaml:"default_security_groups"`
					Subnets               []struct {
						ID            string `yaml:"id"`
						Range         string `yaml:"range"`
						AZ            string `yaml:"az"`
						SecurityGroup string `yaml:"security_group"`
					} `yaml:"subnets"`
					CloudConfigSubnets []struct {
						ID            string `yaml:"id"`
						Range         string `yaml:"range"`
						AZ            string `yaml:"az"`
						SecurityGroup string `yaml:"security_group"`
					} `yaml:"cloud_config_subnets"`
				} `yaml:"aws"`
				BOSH struct {
					Target         string `yaml:"target"`
					Username       string `yaml:"username"`
					Password       string `yaml:"password"`
					DirectorCACert string `yaml:"director_ca_cert"`
					Errand         struct {
						Network struct {
							Name     string `yaml:"name"`
							StaticIP string `yaml:"static_ip"`
							AZ       string `yaml:"az"`
						} `yaml:"network"`
						DefaultPersistentDiskType string `yaml:"default_persistent_disk_type"`
						DefaultVMType             string `yaml:"default_vm_type"`
					} `yaml:"errand"`
				} `yaml:"bosh"`
				Registry struct {
					Username string      `yaml:"username"`
					Password string      `yaml:"password"`
					Host     interface{} `yaml:"host"`
					Port     interface{} `yaml:"port"`
				} `yaml:"registry"`
				ParallelNodes              int    `yaml:"parallel_nodes"`
				ConsulReleaseVersion       string `yaml:"consul_release_version"`
				LatestConsulReleaseVersion string `yaml:"latest_consul_release_version"`
				EnableTurbulenceTests      bool   `yaml:"enable_turbulence_tests"`
			} `yaml:"acceptance_tests"`
		} `yaml:"consul"`
	} `yaml:"properties"`
	Update interface{} `yaml:"update"`
}

func main() {
	output, err := Generate(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(output))
}

func Generate(exampleManifestFilePath string) ([]byte, error) {
	contents, err := ioutil.ReadFile(exampleManifestFilePath)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = yaml.Unmarshal(contents, &manifest)
	if err != nil {
		return nil, err
	}

	manifest.DirectorUUID = os.Getenv("BOSH_DIRECTOR_UUID")
	manifest.Releases[0].Version = os.Getenv("CONSUL_RELEASE_VERSION")
	manifest.Stemcells[0].Version = os.Getenv("STEMCELL_VERSION")
	manifest.InstanceGroups[0].AZs = []string{os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_NETWORK_AZ")}
	manifest.InstanceGroups[0].Networks = []Network{
		{
			Name:      os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_NETWORK_NAME"),
			StaticIPs: []string{os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_NETWORK_STATIC_IP")},
		},
	}
	manifest.Properties.Consul.AcceptanceTests.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	manifest.Properties.Consul.AcceptanceTests.AWS.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	manifest.Properties.Consul.AcceptanceTests.AWS.Region = os.Getenv("AWS_REGION")
	manifest.Properties.Consul.AcceptanceTests.AWS.DefaultSecurityGroups = []string{os.Getenv("AWS_SECURITY_GROUP_NAME")}
	manifest.Properties.Consul.AcceptanceTests.AWS.DefaultKeyName = os.Getenv("AWS_DEFAULT_KEY_NAME")

	manifest.Properties.Consul.AcceptanceTests.BOSH.Target = os.Getenv("BOSH_DIRECTOR")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Username = os.Getenv("BOSH_USER")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Password = os.Getenv("BOSH_PASSWORD")
	manifest.Properties.Consul.AcceptanceTests.BOSH.DirectorCACert = os.Getenv("BOSH_DIRECTOR_CA_CERT")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Errand.Network.AZ = os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_NETWORK_AZ")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Errand.Network.Name = os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_NETWORK_NAME")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Errand.Network.StaticIP = os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_NETWORK_STATIC_IP")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Errand.DefaultVMType = os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_DEFAULT_VM_TYPE")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Errand.DefaultPersistentDiskType = os.Getenv("BOSH_ERRAND_CLOUD_CONFIG_DEFAULT_PERSISTENT_DISK_TYPE")

	manifest.Properties.Consul.AcceptanceTests.Registry.Host = os.Getenv("REGISTRY_HOST")
	manifest.Properties.Consul.AcceptanceTests.Registry.Username = os.Getenv("REGISTRY_USERNAME")
	manifest.Properties.Consul.AcceptanceTests.Registry.Password = os.Getenv("REGISTRY_PASSWORD")

	manifest.Properties.Consul.AcceptanceTests.ConsulReleaseVersion = os.Getenv("CONSUL_RELEASE_VERSION")
	manifest.Properties.Consul.AcceptanceTests.LatestConsulReleaseVersion = os.Getenv("LATEST_CONSUL_RELEASE_VERSION")
	manifest.Properties.Consul.AcceptanceTests.EnableTurbulenceTests = (os.Getenv("ENABLE_TURBULENCE_TESTS") == "true")

	parallelNodes, err := strconv.Atoi(os.Getenv("PARALLEL_NODES"))
	if err != nil {
		return nil, err
	}
	manifest.Properties.Consul.AcceptanceTests.ParallelNodes = parallelNodes

	if err := json.Unmarshal([]byte(os.Getenv("AWS_SUBNETS")), &manifest.Properties.Consul.AcceptanceTests.AWS.Subnets); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(os.Getenv("AWS_CLOUD_CONFIG_SUBNETS")), &manifest.Properties.Consul.AcceptanceTests.AWS.CloudConfigSubnets); err != nil {
		return nil, err
	}

	contents, err = yaml.Marshal(manifest)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
