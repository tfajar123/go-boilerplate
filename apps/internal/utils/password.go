package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024 // 64MB
	iterations  = 3
	parallelism = 2
	keyLen      = 32
	saltLen     = 16
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s.%s", b64Salt, b64Hash), nil
}

func VerifyPassword(hashEncoded, password string) bool {
	parts := strings.Split(hashEncoded, ".")
	if len(parts) != 2 {
		return false
	}

	saltB64 := parts[0]
	hashB64 := parts[1]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(hashB64)
	if err != nil {
		return false
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLen,
	)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}
