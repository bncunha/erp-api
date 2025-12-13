package domain

import "strings"

type Role string

const (
	UserRoleReseller Role = "RESELLER"
	UserRoleAdmin    Role = "ADMIN"
)

type User struct {
	Id          int64
	Username    string
	Name        string
	PhoneNumber *string
	Password    string
	Role        string
	TenantId    int64
	Email       string
}

type CreateUserParams struct {
	Username    string
	Name        string
	PhoneNumber *string
	Role        string
	Email       string
}

func NewUser(params CreateUserParams) User {
	var phone *string
	if params.PhoneNumber != nil {
		trimmed := strings.TrimSpace(*params.PhoneNumber)
		if trimmed != "" {
			phone = &trimmed
		}
	}

	return User{
		Username:    params.Username,
		Name:        params.Name,
		PhoneNumber: phone,
		Role:        params.Role,
		Email:       params.Email,
	}
}
