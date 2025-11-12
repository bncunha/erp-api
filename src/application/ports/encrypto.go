package ports

type Encrypto interface {
	Encrypt(text string) (string, error)
	Compare(hash string, text string) (bool, error)
}