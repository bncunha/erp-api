package ports

import "testing"

type fakeEncrypto struct{}

func (fakeEncrypto) Encrypt(text string) (string, error) { return "enc-" + text, nil }
func (fakeEncrypto) Compare(hash string, text string) (bool, error) {
	return hash == "enc-"+text, nil
}

type fakeEmailPort struct {
	last struct {
		to      string
		subject string
		body    string
	}
}

func (f *fakeEmailPort) Send(to string, subject string, body string) error {
	f.last = struct {
		to      string
		subject string
		body    string
	}{to: to, subject: subject, body: body}
	return nil
}

func TestNewPorts(t *testing.T) {
	encrypto := fakeEncrypto{}
	emailPort := &fakeEmailPort{}
	ports := NewPorts(encrypto, emailPort)
	if ports.Encrypto == nil {
		t.Fatalf("expected encrypto implementation to be set")
	}
	ok, err := ports.Encrypto.Compare("enc-secret", "secret")
	if err != nil || !ok {
		t.Fatalf("expected encrypto implementation to be wired correctly")
	}
	if ports.EmailPort == nil {
		t.Fatalf("expected email port to be set")
	}
	if err := ports.EmailPort.Send("to@example.com", "subject", "body"); err != nil {
		t.Fatalf("expected email port to be callable")
	}
}
