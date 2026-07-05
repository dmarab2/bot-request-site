package auth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	doTheyMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return doTheyMatch, nil

}

func GenerateClaimPassword() (string, error) {
	// 256 bits of entropy is probably overkill but better safe than sorry
	byteArray := make([]byte, 32)
	if _, err := rand.Read(byteArray); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(byteArray), nil
}
