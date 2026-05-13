package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type middlewareStore struct {
	record *apikey.ApiKeyRecord
	err    error
}

func (s middlewareStore) Create(ctx context.Context, record apikey.ApiKeyRecord) (apikey.ApiKeyRecord, error) {
	return record, nil
}
func (s middlewareStore) FindByPublicKey(ctx context.Context, publicKey string) (*apikey.ApiKeyRecord, error) {
	return s.record, s.err
}
func (s middlewareStore) PublicKeyExists(ctx context.Context, publicKey string) (bool, error) {
	return false, nil
}
func (s middlewareStore) MarkUsed(ctx context.Context, id int64, at time.Time, ip *string, userAgent *string) error {
	return nil
}
func (s middlewareStore) Revoke(ctx context.Context, id int64, at time.Time) error { return nil }

func TestAPIKeyAuthMissing(t *testing.T) {
	service := newMiddlewareService(t, middlewareStore{})
	handler := APIKeyAuth(service, "local", "demo:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestAPIKeyAuthValid(t *testing.T) {
	pepper := []byte("01234567890123456789012345678901")
	hash, _ := apikey.HashSecret("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", pepper)
	service, err := apikey.NewService(middlewareStore{record: &apikey.ApiKeyRecord{
		ID:        1,
		App:       "crm",
		Env:       "local",
		PublicKey: "7F3K9Q2M8N4P6R1T",
		Hash:      hash,
		Scopes:    []string{"demo:read"},
		CreatedAt: time.Now().UTC(),
	}}, apikey.Options{Pepper: pepper}, apikey.NoopUsageRecorder{})
	if err != nil {
		t.Fatal(err)
	}
	handler := APIKeyAuth(service, "local", "demo:read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok || principal == nil {
			t.Fatal("principal missing")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	req.Header.Set("Authorization", "Bearer ak_crm_7F3K9Q2M8N4P6R1T_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d", rec.Code)
	}
}

func newMiddlewareService(t *testing.T, store middlewareStore) *apikey.Service {
	t.Helper()
	svc, err := apikey.NewService(store, apikey.Options{Pepper: []byte("01234567890123456789012345678901")}, apikey.NoopUsageRecorder{})
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestWriteValidationFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteValidationFailure(rec, apikey.FailureStoreUnavailable)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d", rec.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
}
