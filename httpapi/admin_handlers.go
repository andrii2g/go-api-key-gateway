package httpapi

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type AdminHandlers struct {
	Service    *apikey.Service
	AdminToken string
	Now        func() time.Time
}

func (h AdminHandlers) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	log.Printf("admin create request started remote_addr=%s", r.RemoteAddr)
	if !validAdminToken(r.Header.Get("X-Admin-Token"), h.AdminToken) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var req struct {
		App       string     `json:"app"`
		Env       string     `json:"env"`
		Name      *string    `json:"name"`
		CreatedBy *string    `json:"created_by"`
		Scopes    []string   `json:"scopes"`
		ExpiresAt *time.Time `json:"expires_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("admin create decode failed: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("admin create decoded app=%q env=%q scopes=%d", req.App, req.Env, len(req.Scopes))
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	result, err := h.Service.Create(ctx, apikey.CreateRequest{
		App:       req.App,
		Env:       req.Env,
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
		Scopes:    req.Scopes,
		ExpiresAt: req.ExpiresAt,
	})
	if err != nil {
		log.Printf("admin create failed: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("admin create success id=%d public_key=%s", result.ID, result.PublicKey)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h AdminHandlers) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	if !validAdminToken(r.Header.Get("X-Admin-Token"), h.AdminToken) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	idValue := r.PathValue("id")
	if idValue == "" {
		idValue = pathTailID(r.URL.Path)
	}
	id, err := strconv.ParseInt(idValue, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	if h.Now != nil {
		now = h.Now().UTC()
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if err := h.Service.Revoke(ctx, id, now); err != nil {
		log.Printf("admin revoke failed: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"revoked": true})
}

func validAdminToken(actual, expected string) bool {
	if actual == "" || expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}

func pathTailID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2]
}
