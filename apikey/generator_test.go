package apikey

import "testing"

func TestGeneratePublicKey(t *testing.T) {
	keys := make(map[string]struct{})
	for range 1000 {
		key, err := GeneratePublicKey()
		if err != nil {
			t.Fatalf("GeneratePublicKey error: %v", err)
		}
		if len(key) != PublicKeyLength {
			t.Fatalf("len = %d", len(key))
		}
		if !isValidPublicKey(key) {
			t.Fatalf("invalid alphabet: %q", key)
		}
		keys[key] = struct{}{}
	}
	if len(keys) == 1 {
		t.Fatal("all generated public keys were identical")
	}
}

func TestGenerateSecretLengths(t *testing.T) {
	tests := []struct {
		bytes int
		want  int
	}{
		{bytes: 32, want: 43},
		{bytes: 48, want: 64},
		{bytes: 64, want: 86},
	}
	for _, tt := range tests {
		secret, err := GenerateSecret(tt.bytes)
		if err != nil {
			t.Fatalf("GenerateSecret(%d) error: %v", tt.bytes, err)
		}
		if len(secret) != tt.want {
			t.Fatalf("GenerateSecret(%d) len=%d want=%d", tt.bytes, len(secret), tt.want)
		}
	}
}

func TestGenerateSecretRejectsSmallValue(t *testing.T) {
	if _, err := GenerateSecret(31); err != ErrInvalidSecretBytes {
		t.Fatalf("err = %v", err)
	}
}

func TestBuildFullKey(t *testing.T) {
	got := BuildFullKey("crm", "7F3K9Q2M8N4P6R1T", "secret")
	want := "ak_crm_7F3K9Q2M8N4P6R1T_secret"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
