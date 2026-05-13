package apikey

import "testing"

func TestHashSecret(t *testing.T) {
	pepper := make([]byte, 32)
	hash, err := HashSecret("secret", pepper)
	if err != nil {
		t.Fatalf("HashSecret error: %v", err)
	}
	if len(hash) != 64 || !isLowerHex(hash) {
		t.Fatalf("unexpected hash %q", hash)
	}
}

func TestCompareSecretHash(t *testing.T) {
	pepper := []byte("01234567890123456789012345678901")
	hash, _ := HashSecret("secret", pepper)

	ok, err := CompareSecretHash("secret", hash, pepper)
	if err != nil || !ok {
		t.Fatalf("expected match ok=%v err=%v", ok, err)
	}

	ok, err = CompareSecretHash("different", hash, pepper)
	if err != nil || ok {
		t.Fatalf("expected mismatch ok=%v err=%v", ok, err)
	}

	ok, err = CompareSecretHash("secret", "INVALID", pepper)
	if err != nil || ok {
		t.Fatalf("expected malformed stored hash to fail safely ok=%v err=%v", ok, err)
	}
}
