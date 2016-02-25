package fakes

import "github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"

type AWS struct {
	SessionCall struct {
		Returns struct {
			Session clients.Session
			Error   error
		}
	}
	FetchSubnetsCall struct {
		Receives struct {
			SubnetIds []string
		}
		Returns struct {
			Subnets []clients.Subnet
			Error   error
		}
	}
}

func (a *AWS) Session() (clients.Session, error) {
	return a.SessionCall.Returns.Session, a.SessionCall.Returns.Error
}

func (a *AWS) FetchSubnets(session clients.Session, subnetIds []string) ([]clients.Subnet, error) {
	a.FetchSubnetsCall.Receives.SubnetIds = subnetIds
	return a.FetchSubnetsCall.Returns.Subnets, a.FetchSubnetsCall.Returns.Error
}
