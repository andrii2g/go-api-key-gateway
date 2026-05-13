package apikey

import (
	"regexp"
	"slices"
	"strings"
)

var scopePattern = regexp.MustCompile(`^[a-z0-9:._*-]{1,100}$`)

func ValidateScopes(scopes []string) error {
	_, err := NormalizeScopes(scopes)
	return err
}

func NormalizeScopes(scopes []string) ([]string, error) {
	if len(scopes) == 0 {
		return []string{}, nil
	}

	seen := make(map[string]struct{}, len(scopes))
	out := make([]string, 0, len(scopes))
	for _, raw := range scopes {
		scope := strings.TrimSpace(raw)
		if !scopePattern.MatchString(scope) {
			return nil, ErrInvalidScope
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	slices.Sort(out)
	return out, nil
}

func HasRequiredScopes(actual []string, required []string) bool {
	if len(required) == 0 {
		return true
	}
	actualSet := make(map[string]struct{}, len(actual))
	for _, scope := range actual {
		actualSet[scope] = struct{}{}
	}
	if _, ok := actualSet["*"]; ok {
		return true
	}
	for _, scope := range required {
		if _, ok := actualSet[scope]; !ok {
			return false
		}
	}
	return true
}
