package ports

import "github.com/bncunha/erp-api/src/domain"

type Ports struct {
	Encrypto  domain.Encrypto
	EmailPort EmailPort
}

func NewPorts(
	encrypto domain.Encrypto,
	emailPort EmailPort,
) *Ports {
	return &Ports{
		Encrypto:  encrypto,
		EmailPort: emailPort,
	}
}
