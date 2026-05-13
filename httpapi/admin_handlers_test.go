package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type adminStore struct {
	record  apikey.ApiKeyRecord
	revoked bool
}

func (s *adminStore) Create(ctx context.Context, record apikey.ApiKeyRecord) (apikey.ApiKeyRecord, error) {
	record.ID = 1
	record.CreatedAt = time.Now().UTC()
	s.record = record
	return record, nil
}
func (s *adminStore) FindByPublicKey(ctx context.Context, publicKey string) (*apikey.ApiKeyRecord, error) {
	return nil, nil
}
func (s *adminStore) PublicKeyExists(ctx context.Context, publicKey string) (bool, error) {
	return false, nil
}
func (s *adminStore) MarkUsed(ctx context.Context, id int64, at time.Time, ip *string, userAgent *string) error {
	return nil
}
func (s *adminStore) Revoke(ctx context.Context, id int64, at time.Time) error {
	s.revoked = true
	return nil
}

func TestAdminCreateRequiresToken(t *testing.T) {
	h := newAdminHandlers(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/api-keys", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	h.CreateAPIKey(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestAdminCreateReturnsAPIKey(t *testing.T) {
	h := newAdminHandlers(t)
	body := []byte(`{"app":"crm","env":"local","scopes":["demo:read"]}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/api-keys", bytes.NewReader(body))
	req.Header.Set("X-Admin-Token", "token")
	rec := httptest.NewRecorder()
	h.CreateAPIKey(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["api_key"] == "" {
		t.Fatal("api_key missing")
	}
}

func TestAdminRevokeRequiresToken(t *testing.T) {
	h := newAdminHandlers(t)
	req := httptest.NewRequest(http.MethodPost, "/admin/api-keys/1/revoke", nil)
	rec := httptest.NewRecorder()
	h.RevokeAPIKey(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rec.Code)
	}
}

func newAdminHandlers(t *testing.T) AdminHandlers {
	t.Helper()
	store := &adminStore{}
	service, err := apikey.NewService(store, apikey.Options{Pepper: []byte("01234567890123456789012345678901")}, apikey.NoopUsageRecorder{})
	if err != nil {
		t.Fatal(err)
	}
	return AdminHandlers{Service: service, AdminToken: "token"}
}
