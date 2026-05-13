package sqlstore

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type ApiKeyRow struct {
	ID                int64
	App               string
	Env               string
	PublicKey         string
	Hash              string
	Scopes            string
	Name              sql.NullString
	CreatedBy         sql.NullString
	CreatedAt         time.Time
	ExpiresAt         sql.NullTime
	RevokedAt         sql.NullTime
	LastUsedAt        sql.NullTime
	LastUsedIP        sql.NullString
	LastUsedUserAgent sql.NullString
}

type InsertParams struct {
	App       string
	Env       string
	PublicKey string
	Hash      string
	Scopes    string
	Name      any
	CreatedBy any
	CreatedAt time.Time
	ExpiresAt any
}

func RowToRecord(row ApiKeyRow) (apikey.ApiKeyRecord, error) {
	var scopes []string
	if err := json.Unmarshal([]byte(row.Scopes), &scopes); err != nil {
		return apikey.ApiKeyRecord{}, err
	}
	return apikey.ApiKeyRecord{
		ID:                row.ID,
		App:               row.App,
		Env:               row.Env,
		PublicKey:         row.PublicKey,
		Hash:              row.Hash,
		Scopes:            scopes,
		Name:              nullStringPtr(row.Name),
		CreatedBy:         nullStringPtr(row.CreatedBy),
		CreatedAt:         row.CreatedAt.UTC(),
		ExpiresAt:         nullTimePtr(row.ExpiresAt),
		RevokedAt:         nullTimePtr(row.RevokedAt),
		LastUsedAt:        nullTimePtr(row.LastUsedAt),
		LastUsedIP:        nullStringPtr(row.LastUsedIP),
		LastUsedUserAgent: nullStringPtr(row.LastUsedUserAgent),
	}, nil
}

func RecordToInsertParams(record apikey.ApiKeyRecord) (InsertParams, error) {
	scopes, err := json.Marshal(record.Scopes)
	if err != nil {
		return InsertParams{}, err
	}
	return InsertParams{
		App:       record.App,
		Env:       record.Env,
		PublicKey: record.PublicKey,
		Hash:      record.Hash,
		Scopes:    string(scopes),
		Name:      stringPtrValue(record.Name),
		CreatedBy: stringPtrValue(record.CreatedBy),
		CreatedAt: record.CreatedAt.UTC(),
		ExpiresAt: timePtrValue(record.ExpiresAt),
	}, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	v := value.String
	return &v
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time.UTC()
	return &t
}

func stringPtrValue(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func timePtrValue(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.UTC()
}
