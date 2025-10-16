package ksuid

import (
	"crypto/rand"
	"encoding/base64"
)

type KSUID struct {
	value string
}

func New() KSUID {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return KSUID{value: base64.RawURLEncoding.EncodeToString(b)}
}

func (k KSUID) String() string {
	return k.value
}
