package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func LoadPepper(base64Value, filePath string) ([]byte, error) {
	switch {
	case strings.TrimSpace(base64Value) != "":
		return decodePepper(strings.TrimSpace(base64Value))
	case strings.TrimSpace(filePath) != "":
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return decodePepper(strings.TrimSpace(string(content)))
	default:
		return nil, fmt.Errorf("pepper source is required")
	}
}

func decodePepper(value string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	if len(decoded) < 32 {
		return nil, fmt.Errorf("decoded pepper must be at least 32 bytes")
	}
	return decoded, nil
}
