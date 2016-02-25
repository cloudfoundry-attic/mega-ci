package clients

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Subnet struct {
	SubnetID  string
	CIDRBlock string
}

type AWS struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type Session interface {
	DescribeSubnets(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
}

func NewAWS(AccessKeyID string, SecretAccessKey string, Region string) AWS {
	return AWS{
		AccessKeyID:     AccessKeyID,
		SecretAccessKey: SecretAccessKey,
		Region:          Region,
	}
}

func (a AWS) FetchSubnets(session Session, subnetIds []string) ([]Subnet, error) {
	var awsSubnetIds []*string
	for _, id := range subnetIds {
		awsSubnetIds = append(awsSubnetIds, aws.String(id))
	}

	params := &ec2.DescribeSubnetsInput{
		SubnetIds: awsSubnetIds,
	}

	resp, err := session.DescribeSubnets(params)
	if err != nil {
		return []Subnet{}, err
	}

	var subnets []Subnet
	for _, subnet := range resp.Subnets {
		subnets = append(subnets, Subnet{
			SubnetID:  *subnet.SubnetId,
			CIDRBlock: *subnet.CidrBlock,
		})
	}

	return subnets, nil
}

func (a AWS) Session() (Session, error) {
	if a.AccessKeyID == "" {
		return nil, errors.New("aws access key id must be provided")
	}

	if a.SecretAccessKey == "" {
		return nil, errors.New("aws secret access key must be provided")
	}

	if a.Region == "" {
		return nil, errors.New("aws region must be provided")
	}

	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(a.AccessKeyID, a.SecretAccessKey, ""),
		Region:      aws.String(a.Region),
	}

	return ec2.New(session.New(awsConfig)), nil
}
