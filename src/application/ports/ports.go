package ports

type Ports struct {
	Encrypto Encrypto
}

func NewPorts(
	encrypto Encrypto,
) *Ports {
	return &Ports{
		Encrypto: encrypto,
	}
}