package model

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

type UpdateRequest struct {
	Old          string `json:"old"`
	New          string `json:"new"`
	SessionToken string `json:"-"`
}

func StringToSaltAndHash(s string) ([]byte, []byte, error) {
	salt := make([]byte, 2056)
	_, err := rand.Read(salt)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	hash := argon2.IDKey([]byte(s), salt, 1, 64*1024, 4, 32)

	return salt, hash, nil
}
