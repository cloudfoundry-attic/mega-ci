package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

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
	manifest.Compilation.CloudProperties.AvailibilityZone = os.Getenv("AWS_AVAILIBILITY_ZONE")
	manifest.Networks[0].Subnets[0].CloudProperties.Subnet = os.Getenv("AWS_SUBNET_ID")
	manifest.Properties.Consul.AcceptanceTests.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	manifest.Properties.Consul.AcceptanceTests.AWS.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	manifest.Properties.Consul.AcceptanceTests.AWS.Region = os.Getenv("AWS_REGION")
	manifest.Properties.Consul.AcceptanceTests.AWS.Subnet = os.Getenv("AWS_SUBNET_ID")
	manifest.Properties.Consul.AcceptanceTests.AWS.DefaultSecurityGroups = []string{os.Getenv("AWS_SECURITY_GROUP_NAME")}
	manifest.Properties.Consul.AcceptanceTests.BOSH.Target = os.Getenv("BOSH_TARGET")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Username = os.Getenv("BOSH_USERNAME")
	manifest.Properties.Consul.AcceptanceTests.BOSH.Password = os.Getenv("BOSH_PASSWORD")
	manifest.Properties.Consul.AcceptanceTests.BOSH.DirectorCACert = os.Getenv("BOSH_DIRECTOR_CA_CERT")
	manifest.Properties.Consul.AcceptanceTests.Registry.Username = os.Getenv("REGISTRY_USERNAME")
	manifest.Properties.Consul.AcceptanceTests.Registry.Password = os.Getenv("REGISTRY_PASSWORD")
	manifest.ResourcePools[0].CloudProperties.AvailibilityZone = os.Getenv("AWS_AVAILIBILITY_ZONE")

	contents, err = yaml.Marshal(manifest)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

type Manifest struct {
	Name         interface{} `yaml:"name"`
	DirectorUUID string      `yaml:"director_uuid"`
	Releases     interface{} `yaml:"releases"`
	Jobs         interface{} `yaml:"jobs"`
	Compilation  struct {
		Workers             interface{} `yaml:"workers"`
		Network             interface{} `yaml:"network"`
		ReuseCompilationVMs interface{} `yaml:"reuse_compilation_vms"`
		CloudProperties     struct {
			AvailibilityZone string      `yaml:"availability_zone"`
			EphemeralDisk    interface{} `yaml:"ephemeral_disk"`
			InstanceType     interface{} `yaml:"instance_type"`
		} `yaml:"cloud_properties"`
	} `yaml:"compilation"`
	Networks []struct {
		Name    interface{} `yaml:"name"`
		Type    interface{} `yaml:"type"`
		Subnets []struct {
			Range           interface{} `yaml:"range"`
			Gateway         interface{} `yaml:"gateway"`
			Static          interface{} `yaml:"static"`
			Reserved        interface{} `yaml:"reserved"`
			CloudProperties struct {
				Subnet string `yaml:"subnet"`
			} `yaml:"cloud_properties"`
		} `yaml:"subnets"`
	} `yaml:"networks"`
	Properties struct {
		Consul struct {
			AcceptanceTests struct {
				AWS struct {
					AccessKeyID           string      `yaml:"access_key_id"`
					SecretAccessKey       string      `yaml:"secret_access_key"`
					Region                string      `yaml:"region"`
					DefaultKeyName        interface{} `yaml:"default_key_name"`
					DefaultSecurityGroups []string    `yaml:"default_security_groups"`
					Subnet                string      `yaml:"subnet"`
				} `yaml:"aws"`
				BOSH struct {
					Target         string `yaml:"target"`
					Username       string `yaml:"username"`
					Password       string `yaml:"password"`
					DirectorCACert string `yaml:"director_ca_cert"`
				} `yaml:"bosh"`
				Registry struct {
					Username string      `yaml:"username"`
					Password string      `yaml:"password"`
					Host     interface{} `yaml:"host"`
					Port     interface{} `yaml:"port"`
				} `yaml:"registry"`
			} `yaml:"acceptance_tests"`
		} `yaml:"consul"`
	} `yaml:"properties"`
	ResourcePools []struct {
		Name            interface{} `yaml:"name"`
		Network         interface{} `yaml:"network"`
		Stemcell        interface{} `yaml:"stemcell"`
		CloudProperties struct {
			AvailibilityZone string      `yaml:"availability_zone"`
			EphemeralDisk    interface{} `yaml:"ephemeral_disk"`
			InstanceType     interface{} `yaml:"instance_type"`
		} `yaml:"cloud_properties"`
	} `yaml:"resource_pools"`
	Update interface{} `yaml:"update"`
}
