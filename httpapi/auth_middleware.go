package httpapi

import (
	"net"
	"net/http"
	"strings"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

func APIKeyAuth(service *apikey.Service, env string, requiredScopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawKey := ""
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				rawKey = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			} else if authHeader != "" {
				WriteValidationFailure(w, apikey.FailureMalformed)
				return
			}

			ip := remoteIP(r)
			userAgent := r.UserAgent()
			result := service.Validate(r.Context(), apikey.ValidationRequest{
				RawKey:         rawKey,
				Environment:    env,
				RequiredScopes: requiredScopes,
				IP:             ip,
				UserAgent:      stringPointer(userAgent),
			})
			if !result.OK {
				WriteValidationFailure(w, result.Reason)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), result.Principal)))
		})
	}
}

func remoteIP(r *http.Request) *string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if r.RemoteAddr == "" {
			return nil
		}
		return &r.RemoteAddr
	}
	return &host
}

func stringPointer(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}
