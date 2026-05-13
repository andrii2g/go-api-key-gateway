package apikey

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GeneratePublicKey() (string, error) {
	buf := make([]byte, PublicKeyLength)
	for i := 0; i < PublicKeyLength; i++ {
		index, err := randomAlphabetIndex(len(PublicKeyAlphabet))
		if err != nil {
			return "", err
		}
		buf[i] = PublicKeyAlphabet[index]
	}
	return string(buf), nil
}

func GenerateSecret(secretBytes int) (string, error) {
	if secretBytes < MinSecretBytes {
		return "", ErrInvalidSecretBytes
	}
	b := make([]byte, secretBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func BuildFullKey(app, publicKey, secret string) string {
	return fmt.Sprintf("%s_%s_%s_%s", KeyPrefix, app, publicKey, secret)
}

func randomAlphabetIndex(size int) (int, error) {
	if size <= 0 || size > 256 {
		return 0, ErrInvalidSecretBytes
	}

	max := 256 - (256 % size)
	var b [1]byte
	for {
		if _, err := rand.Read(b[:]); err != nil {
			return 0, err
		}
		value := int(b[0])
		if value < max {
			return value % size, nil
		}
	}
}
