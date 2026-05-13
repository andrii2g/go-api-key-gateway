package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPepperFromBase64(t *testing.T) {
	value := base64.StdEncoding.EncodeToString([]byte("01234567890123456789012345678901"))
	pepper, err := LoadPepper(value, "")
	if err != nil {
		t.Fatalf("LoadPepper error: %v", err)
	}
	if len(pepper) != 32 {
		t.Fatalf("len=%d", len(pepper))
	}
}

func TestLoadPepperFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pepper")
	value := base64.StdEncoding.EncodeToString([]byte("01234567890123456789012345678901"))
	if err := os.WriteFile(path, []byte(value+"\n"), 0o600); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	pepper, err := LoadPepper("", path)
	if err != nil {
		t.Fatalf("LoadPepper error: %v", err)
	}
	if len(pepper) != 32 {
		t.Fatalf("len=%d", len(pepper))
	}
}
