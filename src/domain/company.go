package domain

type Company struct {
	Id        int64
	Name      string
	LegalName string
	Cnpj      string
	Cpf       string
	Cellphone string
	Address   *Address
}
