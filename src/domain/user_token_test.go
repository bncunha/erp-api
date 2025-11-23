package domain

import (
	"testing"
	"time"
)

type stubTokenEncrypto struct {
	encryptErr error
	lastText   string
}

func (s *stubTokenEncrypto) Encrypt(text string) (string, error) {
	s.lastText = text
	if s.encryptErr != nil {
		return "", s.encryptErr
	}
	return "hashed-" + text, nil
}

func (s *stubTokenEncrypto) Compare(hash string, text string) (bool, error) {
	return hash == "hashed-"+text, nil
}

func TestNewUserToken(t *testing.T) {
	encrypto := &stubTokenEncrypto{}
	token := NewUserToken(CreateUserTokenParams{
		User:      User{Id: 1},
		CreatedBy: User{Id: 2},
		Type:      UserTokenTypeInvite,
	}, encrypto)

	if token.Code == "" || token.Uuid == "" {
		t.Fatalf("expected code and uuid to be generated: %+v", token)
	}
	if token.CodeHash != "hashed-"+token.Code {
		t.Fatalf("expected code hash to be generated")
	}
	if time.Until(token.ExpiresAt) <= 0 {
		t.Fatalf("expected expiration to be in the future")
	}
	if encrypto.lastText == "" {
		t.Fatalf("expected encrypto to be invoked")
	}
}

func TestUserTokenIsValid(t *testing.T) {
	encrypto := &stubTokenEncrypto{}
	token := UserToken{
		CodeHash:  "hashed-secret",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	ok, err := token.IsValid(encrypto, "secret")
	if err != nil || !ok {
		t.Fatalf("expected token to be valid")
	}

	token.UsedAt = ptrTime(time.Now())
	if ok, _ := token.IsValid(encrypto, "secret"); ok {
		t.Fatalf("expected used token to be invalid")
	}
}

func TestUserTokenSetUsedAt(t *testing.T) {
	token := UserToken{}
	token.SetUsedAt()
	if token.UsedAt == nil {
		t.Fatalf("expected used at to be set")
	}
}

func TestGenerateAndSetCodeHashError(t *testing.T) {
	encrypto := &stubTokenEncrypto{encryptErr: assertErr{}}
	token := UserToken{}
	if err := token.generateAndSetCodeHash(encrypto); err == nil {
		t.Fatalf("expected error when encrypt fails")
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "fail" }

func ptrTime(t time.Time) *time.Time {
	return &t
}
