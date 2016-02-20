package fakes

type Arguments struct {
	Manifest     string
	BoshDirector string
	BoshUser     string
	BoshPassword string
}

type BOSH struct {
	DeployCalls struct {
		Receives []Arguments
		Returns  struct {
			Error error
		}
	}
	StatusCall struct {
		Returns struct {
			UUID string
		}
	}
}

func (b *BOSH) Deploy(manifest string, boshDirector string, boshUser string, boshPassword string) error {
	b.DeployCalls.Receives = append(b.DeployCalls.Receives, Arguments{
		Manifest:     manifest,
		BoshDirector: boshDirector,
		BoshUser:     boshUser,
		BoshPassword: boshPassword,
	})

	return b.DeployCalls.Returns.Error
}

func (b *BOSH) Status(boshDirector string, boshUser string, boshPassword string) (string, error) {
	return b.StatusCall.Returns.UUID, nil
}
