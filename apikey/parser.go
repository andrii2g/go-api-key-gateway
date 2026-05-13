package apikey

import (
	"encoding/base64"
	"strings"
	"unicode"
)

func Parse(raw string) (ParsedKey, ValidationFailureReason) {
	trimmed := strings.TrimFunc(raw, func(r rune) bool {
		return r <= unicode.MaxASCII && unicode.IsSpace(r)
	})
	if trimmed == "" {
		return ParsedKey{}, FailureMissing
	}

	parts := strings.SplitN(trimmed, "_", 4)
	if len(parts) != 4 {
		return ParsedKey{}, FailureMalformed
	}
	if parts[0] != KeyPrefix {
		return ParsedKey{}, FailureMalformed
	}

	app := parts[1]
	publicKey := parts[2]
	secret := parts[3]

	if !isValidApp(app) {
		return ParsedKey{}, FailureMalformed
	}
	if !isValidPublicKey(publicKey) {
		return ParsedKey{}, FailureMalformed
	}
	if secret == "" || !isBase64URLNoPadding(secret) {
		return ParsedKey{}, FailureMalformed
	}

	decoded, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil || len(decoded) < MinSecretBytes {
		return ParsedKey{}, FailureMalformed
	}

	return ParsedKey{
		App:       app,
		PublicKey: publicKey,
		Secret:    secret,
	}, FailureNone
}

func isBase64URLNoPadding(value string) bool {
	for _, r := range value {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_':
		default:
			return false
		}
	}
	return true
}
