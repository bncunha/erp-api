package viewmodel

import "github.com/bncunha/erp-api/src/domain"

type UserViewModel struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	Email       string `json:"email"`
}

func ToUserViewModel(user domain.User) UserViewModel {
	return UserViewModel{
		Id:          user.Id,
		Username:    user.Username,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		Email:       user.Email,
	}
}
