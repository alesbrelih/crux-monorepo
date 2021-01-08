package pkg

import "golang.org/x/crypto/bcrypt"

type PasswordUtil interface {
	Hash(password []byte, cost int) ([]byte, error)
	Verify(password []byte, hash []byte) bool
}

type passwordUtil struct{}

func (u *passwordUtil) Hash(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}
func (u *passwordUtil) Verify(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err != nil
}

func NewPasswordUtil() PasswordUtil {
	return &passwordUtil{}
}
