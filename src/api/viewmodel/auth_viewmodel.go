package viewmodel

import "github.com/bncunha/erp-api/src/application/service/output"

type LoginViewModel struct {
	Token string `json:"token"`
}

func ToLoginViewModel(out output.LoginOutput) LoginViewModel {
	return LoginViewModel{
		Token: out.Token,
	}
}