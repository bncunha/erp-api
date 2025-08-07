package domain

type User struct {
	Id          int64
	Username    string
	Name        string
	PhoneNumber string
	Password    string
	Role        string
	TenantId    int64
}