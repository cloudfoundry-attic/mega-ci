package subnetchecker

import (
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/manifests"
)

type SubnetChecker struct {
	awsClient AWSClient
}

type AWSClient interface {
	FetchSubnets(session clients.Session, subnetIds []string) ([]clients.Subnet, error)
	Session() (clients.Session, error)
}

func (s SubnetChecker) CheckSubnets(manifestFilename string) (bool, error) {
	networks, err := manifests.ReadNetworksFromManifest(manifestFilename)
	if err != nil {
		return false, err
	}

	manifestSubnetMap := mapSubnets(subnetsFromManifestNetworks(networks))

	session, err := s.awsClient.Session()
	if err != nil {
		return false, err
	}

	awsSubnets, err := s.awsClient.FetchSubnets(session, subnetIdsFromSubnets(manifestSubnetMap))
	if err != nil {
		return false, err
	}
	awsSubnetMap := mapSubnets(awsSubnets)

	for id, manifestSubnet := range manifestSubnetMap {
		awsSubnet, ok := awsSubnetMap[id]
		if !ok {
			return false, nil
		}
		if awsSubnet.CIDRBlock != manifestSubnet.CIDRBlock {
			return false, nil
		}
	}

	return true, nil
}

func NewSubnetChecker(awsClient AWSClient) SubnetChecker {
	return SubnetChecker{
		awsClient: awsClient,
	}
}

func subnetsFromManifestNetworks(networks []manifests.Network) []clients.Subnet {
	var returnedSubnets []clients.Subnet
	for _, network := range networks {
		for _, subnet := range network.Subnets {
			returnedSubnets = append(
				returnedSubnets,
				clients.Subnet{
					SubnetID:  subnet.CloudProperties.Subnet,
					CIDRBlock: subnet.Range,
				},
			)
		}
	}
	return returnedSubnets
}

func subnetIdsFromSubnets(subnets map[string]clients.Subnet) []string {
	var subnetIds []string
	for subnetId, _ := range subnets {
		subnetIds = append(subnetIds, subnetId)
	}
	return subnetIds
}

func mapSubnets(subnets []clients.Subnet) map[string]clients.Subnet {
	subnetMap := make(map[string]clients.Subnet)
	for _, subnet := range subnets {
		subnetMap[subnet.SubnetID] = subnet
	}
	return subnetMap
}
