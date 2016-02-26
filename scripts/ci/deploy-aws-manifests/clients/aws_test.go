package clients_test

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/clients"
	"github.com/cloudfoundry/mega-ci/scripts/ci/deploy-aws-manifests/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS", func() {
	var (
		awsClient   clients.AWS
		fakeSession *fakes.Session
	)
	BeforeEach(func() {
		awsClient = clients.NewAWS("some-access-key-id", "some-secret-access-key", "some-region", "some-endpoint")
		fakeSession = &fakes.Session{}
	})

	Describe("Session", func() {
		It("returns a new ec2 session", func() {
			session, err := awsClient.Session()
			Expect(err).NotTo(HaveOccurred())

			_, ok := session.(clients.Session)
			Expect(ok).To(BeTrue())

			client, ok := session.(*ec2.EC2)
			Expect(ok).To(BeTrue())

			Expect(client.Config.Credentials).To(Equal(credentials.NewStaticCredentials("some-access-key-id", "some-secret-access-key", "")))
			Expect(client.Config.Region).To(Equal(aws.String("some-region")))
			Expect(client.Config.Endpoint).To(Equal(aws.String("some-endpoint")))
		})

		Context("failure cases", func() {
			It("returns an error if no access key id is provided", func() {
				awsClient := clients.NewAWS("", "some-secret-access-key", "some-region", "some-endpoint")
				_, err := awsClient.Session()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("aws access key id must be provided"))
			})

			It("returns an error if no secret access key is provided", func() {
				awsClient := clients.NewAWS("some-access-key-id", "", "some-region", "some-endpoint")
				_, err := awsClient.Session()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("aws secret access key must be provided"))
			})

			It("returns an error if no region is provided", func() {
				awsClient := clients.NewAWS("some-access-key-id", "some-secret-access-key", "", "some-endpoint")
				_, err := awsClient.Session()
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("aws region must be provided"))
			})
		})
	})

	Describe("FetchSubnets", func() {
		It("returns subnets from aws that match subnetIds passed in", func() {
			fakeSession.DescribeSubnetsCall.Returns.DescribeSubnetsOutput = &ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					&ec2.Subnet{
						SubnetId:  aws.String("subnet-1"),
						CidrBlock: aws.String("some-cidr-block"),
					},
				},
			}

			subnets, err := awsClient.FetchSubnets(fakeSession, []string{"subnet-1"})
			Expect(err).NotTo(HaveOccurred())
			Expect(subnets).To(Equal([]clients.Subnet{
				clients.Subnet{
					SubnetID:  "subnet-1",
					CIDRBlock: "some-cidr-block",
				},
			}))
		})
		Context("failure cases", func() {
			It("returns an error if describe subnets fails", func() {
				fakeSession.DescribeSubnetsCall.Returns.Error = errors.New("DescribeSubnets failed")

				_, err := awsClient.FetchSubnets(fakeSession, []string{"subnet-1"})
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("DescribeSubnets failed"))
			})
		})
	})
})
