package passwordhashing

type IPasswordHashing interface {
	Hash(plaintext string) string
	Compare(hash string, plaintext string) bool
}
