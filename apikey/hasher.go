package apikey

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashSecret(secret string, pepper []byte) (string, error) {
	if len(pepper) < MinSecretBytes {
		return "", ErrInvalidPepper
	}
	mac := hmac.New(sha256.New, pepper)
	_, _ = mac.Write([]byte(secret))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func CompareSecretHash(secret string, expectedHashHex string, pepper []byte) (bool, error) {
	if len(expectedHashHex) != 64 || !isLowerHex(expectedHashHex) {
		return false, nil
	}
	actual, err := HashSecret(secret, pepper)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(actual), []byte(expectedHashHex)), nil
}

func isLowerHex(value string) bool {
	for _, r := range value {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}
