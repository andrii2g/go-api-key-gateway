package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

func NewServer(service *apikey.Service, env string, adminToken string) http.Handler {
	mux := http.NewServeMux()
	admin := AdminHandlers{
		Service:    service,
		AdminToken: adminToken,
	}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /admin/api-keys", admin.CreateAPIKey)
	mux.HandleFunc("POST /admin/api-keys/{id}/revoke", admin.RevokeAPIKey)
	mux.Handle("GET /v1/ping", APIKeyAuth(service, env, "demo:read")(http.HandlerFunc(Ping)))
	mux.Handle("POST /v1/messages", APIKeyAuth(service, env, "messages:write")(http.HandlerFunc(Messages)))
	return mux
}
