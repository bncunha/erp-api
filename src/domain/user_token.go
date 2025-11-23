package domain

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type UserTokenType string

const (
	UserTokenTypeInvite    UserTokenType = "INVITE"
	UserTokenTypeResetPass UserTokenType = "RESET_PASSWORD"
)

type UserToken struct {
	Id        int64
	Uuid      string
	User      User
	Type      UserTokenType
	Code      string
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

func NewUserToken(params CreateUserTokenParams, encrypto Encrypto) UserToken {
	ut := UserToken{
		User:      params.User,
		Uuid:      uuid.NewString(),
		Type:      params.Type,
		CreatedBy: params.CreatedBy,
	}
	ut.generateAndSetCodeHash(encrypto)
	ut.setExpiresAt()
	return ut
}

func (ut *UserToken) IsValid(encrypto Encrypto, code string) (bool, error) {
	isValid, err := encrypto.Compare(ut.CodeHash, code)
	if err != nil {
		return false, err
	}
	return ut.ExpiresAt.After(time.Now()) && ut.UsedAt == nil && isValid, nil
}

func (ut *UserToken) SetUsedAt() {
	now := time.Now()
	ut.UsedAt = &now
}

func (ut *UserToken) setExpiresAt() {
	ut.ExpiresAt = time.Now().Add(time.Hour * 1)
}

func (ut *UserToken) generateAndSetCodeHash(encrypto Encrypto) error {
	b := make([]byte, 32)
	rand.Read(b)
	code := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	hash, err := encrypto.Encrypt(code)
	if err != nil {
		return err
	}
	ut.Code = code
	ut.CodeHash = hash
	return nil
}
