package domain

import "time"

type UserTokenType string

const (
	UserTokenTypeInvite    UserTokenType = "invite"
	UserTokenTypeResetPass UserTokenType = "reset_password"
)

type UserToken struct {
	Id           int64
	UserId       User
	Type         UserTokenType
	CodeHash     string
	ExpiresAt    time.Time
	UsedAt       time.Time
	CreatedBy    User
	CreatedAt    time.Time
}
