package ports

type Ports struct {
	Encrypto  Encrypto
	EmailPort EmailPort
}

func NewPorts(
	encrypto Encrypto,
	emailPort EmailPort,
) *Ports {
	return &Ports{
		Encrypto:  encrypto,
		EmailPort: emailPort,
	}
}
