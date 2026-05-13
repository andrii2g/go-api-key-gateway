package apikey

import "time"

type ApiKeyRecord struct {
	ID                int64
	App               string
	Env               string
	PublicKey         string
	Hash              string
	Scopes            []string
	Name              *string
	CreatedBy         *string
	CreatedAt         time.Time
	ExpiresAt         *time.Time
	RevokedAt         *time.Time
	LastUsedAt        *time.Time
	LastUsedIP        *string
	LastUsedUserAgent *string
}

type CreateRequest struct {
	App       string
	Env       string
	Scopes    []string
	Name      *string
	CreatedBy *string
	ExpiresAt *time.Time
}

type CreateResult struct {
	ID        int64      `json:"id"`
	App       string     `json:"app"`
	Env       string     `json:"env"`
	PublicKey string     `json:"public_key"`
	APIKey    string     `json:"api_key"`
	Scopes    []string   `json:"scopes"`
	Name      *string    `json:"name,omitempty"`
	CreatedBy *string    `json:"created_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ParsedKey struct {
	App       string
	PublicKey string
	Secret    string
}

type Principal struct {
	KeyID     int64    `json:"key_id"`
	App       string   `json:"app"`
	Env       string   `json:"env"`
	PublicKey string   `json:"public_key"`
	Scopes    []string `json:"scopes"`
	Name      *string  `json:"name,omitempty"`
}

type ValidationRequest struct {
	RawKey         string
	Environment    string
	RequiredScopes []string
	IP             *string
	UserAgent      *string
}
