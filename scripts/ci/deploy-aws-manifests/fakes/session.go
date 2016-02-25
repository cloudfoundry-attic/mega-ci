package fakes

import "github.com/aws/aws-sdk-go/service/ec2"

type Session struct {
	DescribeSubnetsCall struct {
		Returns struct {
			DescribeSubnetsOutput *ec2.DescribeSubnetsOutput
			Error                 error
		}
	}
}

func (s *Session) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return s.DescribeSubnetsCall.Returns.DescribeSubnetsOutput, s.DescribeSubnetsCall.Returns.Error
}
