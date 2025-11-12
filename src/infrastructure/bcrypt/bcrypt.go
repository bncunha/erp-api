package bcrypt

import (
	"github.com/bncunha/erp-api/src/application/ports"
	_bcrypt "golang.org/x/crypto/bcrypt"
)

type bcrypt struct {
}

func NewBcrypt() ports.Encrypto {
	return &bcrypt{}
}

func (b *bcrypt) Encrypt(text string) (string, error) {
	hash, err := _bcrypt.GenerateFromPassword([]byte(text), _bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *bcrypt) Compare(hash string, text string) (bool, error) {
	err := _bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	if err != nil {
		return false, err
	}
	return true, nil
}