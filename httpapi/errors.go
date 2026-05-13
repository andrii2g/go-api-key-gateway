package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteValidationFailure(w http.ResponseWriter, reason apikey.ValidationFailureReason) {
	status := http.StatusUnauthorized
	body := errorResponse{
		Error:   "invalid_api_key",
		Message: "The API key is missing, malformed, expired, revoked, or invalid.",
	}

	switch reason {
	case apikey.FailureMissing:
		body.Error = "missing_api_key"
		body.Message = "The API key is missing."
	case apikey.FailureScopeDenied:
		status = http.StatusForbidden
		body.Error = "insufficient_scope"
		body.Message = "The API key does not have the required scope."
	case apikey.FailureStoreUnavailable:
		status = http.StatusServiceUnavailable
		body.Error = "api_key_store_unavailable"
		body.Message = "The API key store is unavailable."
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
