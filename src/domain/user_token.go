package domain

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type UserTokenType string

const (
	UserTokenTypeInvite    UserTokenType = "invite"
	UserTokenTypeResetPass UserTokenType = "reset_password"
)

type UserToken struct {
	Id        int64
	User      User
	Type      UserTokenType
	CodeHash  string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedBy User
}

type CreateUserTokenParams struct {
	User      User
	CreatedBy User
	Type      UserTokenType
}

func NewUserToken(params CreateUserTokenParams) UserToken {
	ut := UserToken{
		User:      params.User,
		Type:      params.Type,
		CreatedBy: params.CreatedBy,
	}
	ut.generateAndSetCodeHash()
	ut.setExpiresAt()
	return ut
}

func (ut *UserToken) IsValid() bool {
	return ut.ExpiresAt.After(time.Now()) && ut.UsedAt == nil
}

func (ut *UserToken) setExpiresAt() {
	ut.ExpiresAt = time.Now().Add(time.Hour * 1)
}

func (ut *UserToken) generateAndSetCodeHash() {
	b := make([]byte, 32)
	rand.Read(b)
	ut.CodeHash = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}
