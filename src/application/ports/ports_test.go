package ports

import "testing"

type fakeEncrypto struct{}

func (fakeEncrypto) Encrypt(text string) (string, error) { return "enc-" + text, nil }
func (fakeEncrypto) Compare(hash string, text string) (bool, error) {
	return hash == "enc-"+text, nil
}

func TestNewPorts(t *testing.T) {
	encrypto := fakeEncrypto{}
	ports := NewPorts(encrypto)
	if ports.Encrypto == nil {
		t.Fatalf("expected encrypto implementation to be set")
	}
	ok, err := ports.Encrypto.Compare("enc-secret", "secret")
	if err != nil || !ok {
		t.Fatalf("expected encrypto implementation to be wired correctly")
	}
}
