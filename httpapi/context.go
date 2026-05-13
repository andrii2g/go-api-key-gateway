package httpapi

import (
	"context"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type principalContextKey struct{}

func WithPrincipal(ctx context.Context, principal *apikey.Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (*apikey.Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(*apikey.Principal)
	return principal, ok
}
