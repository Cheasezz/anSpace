package hasher

import (
	"crypto/sha1"
	"fmt"

	"github.com/Cheasezz/anSpace/backend/config"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(cfg config.Hasher) *SHA1Hasher {
	return &SHA1Hasher{salt: cfg.Salt}
}

func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
