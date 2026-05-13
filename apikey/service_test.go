package apikey

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeStore struct {
	createFn          func(ctx context.Context, record ApiKeyRecord) (ApiKeyRecord, error)
	findByPublicKeyFn func(ctx context.Context, publicKey string) (*ApiKeyRecord, error)
	publicKeyExistsFn func(ctx context.Context, publicKey string) (bool, error)
	markUsedFn        func(ctx context.Context, id int64, at time.Time, ip *string, userAgent *string) error
	revokeFn          func(ctx context.Context, id int64, at time.Time) error
}

func (f fakeStore) Create(ctx context.Context, record ApiKeyRecord) (ApiKeyRecord, error) {
	return f.createFn(ctx, record)
}
func (f fakeStore) FindByPublicKey(ctx context.Context, publicKey string) (*ApiKeyRecord, error) {
	return f.findByPublicKeyFn(ctx, publicKey)
}
func (f fakeStore) PublicKeyExists(ctx context.Context, publicKey string) (bool, error) {
	return f.publicKeyExistsFn(ctx, publicKey)
}
func (f fakeStore) MarkUsed(ctx context.Context, id int64, at time.Time, ip *string, userAgent *string) error {
	if f.markUsedFn == nil {
		return nil
	}
	return f.markUsedFn(ctx, id, at, ip, userAgent)
}
func (f fakeStore) Revoke(ctx context.Context, id int64, at time.Time) error {
	if f.revokeFn == nil {
		return nil
	}
	return f.revokeFn(ctx, id, at)
}

func TestServiceCreate(t *testing.T) {
	var saved ApiKeyRecord
	svc, err := NewService(fakeStore{
		createFn: func(ctx context.Context, record ApiKeyRecord) (ApiKeyRecord, error) {
			saved = record
			record.ID = 1
			return record, nil
		},
		publicKeyExistsFn: func(ctx context.Context, publicKey string) (bool, error) {
			return false, nil
		},
	}, Options{Pepper: []byte("01234567890123456789012345678901")}, &MemoryUsageRecorder{})
	if err != nil {
		t.Fatalf("NewService error: %v", err)
	}
	now := time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	result, err := svc.Create(context.Background(), CreateRequest{
		App:       "CRM",
		Env:       "local",
		Scopes:    []string{"messages:write", "demo:read"},
		Name:      stringPtr("CRM local"),
		CreatedBy: stringPtr("admin@example.com"),
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.APIKey == "" {
		t.Fatal("expected api key")
	}
	if saved.Hash == "" || saved.PublicKey == "" || saved.App != "crm" || saved.Env != "local" {
		t.Fatalf("unexpected saved record: %#v", saved)
	}
	if result.CreatedAt != now {
		t.Fatalf("created_at = %v want %v", result.CreatedAt, now)
	}
}

func TestServiceValidate(t *testing.T) {
	recorder := &MemoryUsageRecorder{}
	pepper := []byte("01234567890123456789012345678901")
	hash, _ := HashSecret("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", pepper)
	record := &ApiKeyRecord{
		ID:        1,
		App:       "crm",
		Env:       "local",
		PublicKey: "7F3K9Q2M8N4P6R1T",
		Hash:      hash,
		Scopes:    []string{"demo:read", "messages:write"},
		CreatedAt: time.Now().UTC(),
	}
	svc, err := NewService(fakeStore{
		findByPublicKeyFn: func(ctx context.Context, publicKey string) (*ApiKeyRecord, error) {
			return record, nil
		},
		publicKeyExistsFn: func(ctx context.Context, publicKey string) (bool, error) {
			return false, nil
		},
		createFn: func(ctx context.Context, record ApiKeyRecord) (ApiKeyRecord, error) {
			return record, nil
		},
	}, Options{Pepper: pepper}, recorder)
	if err != nil {
		t.Fatalf("NewService error: %v", err)
	}
	result := svc.Validate(context.Background(), ValidationRequest{
		RawKey:         "ak_crm_7F3K9Q2M8N4P6R1T_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Environment:    "local",
		RequiredScopes: []string{"demo:read"},
		IP:             stringPtr("127.0.0.1"),
		UserAgent:      stringPtr("test"),
	})
	if !result.OK || result.Principal == nil || result.Reason != FailureNone {
		t.Fatalf("unexpected result: %#v", result)
	}
	if len(recorder.Events) != 1 {
		t.Fatalf("events = %d", len(recorder.Events))
	}
}

func TestServiceValidateFailures(t *testing.T) {
	pepper := []byte("01234567890123456789012345678901")
	svc, err := NewService(fakeStore{
		findByPublicKeyFn: func(ctx context.Context, publicKey string) (*ApiKeyRecord, error) {
			return nil, errors.New("db down")
		},
		publicKeyExistsFn: func(ctx context.Context, publicKey string) (bool, error) {
			return false, nil
		},
		createFn: func(ctx context.Context, record ApiKeyRecord) (ApiKeyRecord, error) {
			return record, nil
		},
	}, Options{Pepper: pepper}, NoopUsageRecorder{})
	if err != nil {
		t.Fatalf("NewService error: %v", err)
	}

	if result := svc.Validate(context.Background(), ValidationRequest{}); result.Reason != FailureMissing {
		t.Fatalf("missing reason = %q", result.Reason)
	}
	if result := svc.Validate(context.Background(), ValidationRequest{RawKey: "bad", Environment: "local"}); result.Reason != FailureMalformed {
		t.Fatalf("malformed reason = %q", result.Reason)
	}
	if result := svc.Validate(context.Background(), ValidationRequest{
		RawKey:      "ak_crm_7F3K9Q2M8N4P6R1T_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		Environment: "local",
	}); result.Reason != FailureStoreUnavailable {
		t.Fatalf("store reason = %q", result.Reason)
	}
}

func stringPtr(value string) *string {
	return &value
}
