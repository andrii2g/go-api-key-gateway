package apikey

import (
	"reflect"
	"testing"
)

func TestNormalizeScopes(t *testing.T) {
	got, err := NormalizeScopes([]string{"messages:write", "demo:read", "demo:read"})
	if err != nil {
		t.Fatalf("NormalizeScopes error: %v", err)
	}
	want := []string{"demo:read", "messages:write"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestNormalizeScopesRejectsInvalid(t *testing.T) {
	if _, err := NormalizeScopes([]string{"Demo:Read"}); err != ErrInvalidScope {
		t.Fatalf("err = %v", err)
	}
}

func TestHasRequiredScopes(t *testing.T) {
	if !HasRequiredScopes([]string{"demo:read"}, nil) {
		t.Fatal("empty required should pass")
	}
	if !HasRequiredScopes([]string{"*"}, []string{"messages:write"}) {
		t.Fatal("wildcard should pass")
	}
	if HasRequiredScopes([]string{"demo:read"}, []string{"messages:write"}) {
		t.Fatal("missing scope should fail")
	}
}
