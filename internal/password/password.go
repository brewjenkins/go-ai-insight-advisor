package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidLength = errors.New("password must be between 1 and 72 bytes")

func Hash(raw string) (string, error) {
	if len(raw) == 0 || len(raw) > 72 {
		return "", ErrInvalidLength
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(hash), err
}

func Verify(hash, raw string) bool {
	return len(raw) > 0 && len(raw) <= 72 &&
		bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}
