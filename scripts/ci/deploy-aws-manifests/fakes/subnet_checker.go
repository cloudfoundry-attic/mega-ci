package fakes

type SubnetChecker struct {
	CheckSubnetsCall struct {
		Returns struct {
			Bool  bool
			Error error
		}
	}
}

func (s *SubnetChecker) CheckSubnets(manifestFilename string) (bool, error) {
	return s.CheckSubnetsCall.Returns.Bool, s.CheckSubnetsCall.Returns.Error
}
