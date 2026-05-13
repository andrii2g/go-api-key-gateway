package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/andrii2g/go-api-key-gateway/apikey"
)

type Store struct {
	db         *sql.DB
	statements Statements
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:         db,
		statements: StatementsFor(DialectMySQL),
	}
}

func (s *Store) Create(ctx context.Context, record apikey.ApiKeyRecord) (apikey.ApiKeyRecord, error) {
	params, err := RecordToInsertParams(record)
	if err != nil {
		return apikey.ApiKeyRecord{}, err
	}
	result, err := s.db.ExecContext(ctx, s.statements.InsertAPIKey,
		params.App,
		params.Env,
		params.PublicKey,
		params.Hash,
		params.Scopes,
		params.Name,
		params.CreatedBy,
		params.CreatedAt,
		params.ExpiresAt,
	)
	if err != nil {
		return apikey.ApiKeyRecord{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return apikey.ApiKeyRecord{}, err
	}
	return s.findByID(ctx, id)
}

func (s *Store) FindByPublicKey(ctx context.Context, publicKey string) (*apikey.ApiKeyRecord, error) {
	row, err := s.scanOne(ctx, s.statements.FindByPublicKey, publicKey)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	record, err := RowToRecord(row)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *Store) PublicKeyExists(ctx context.Context, publicKey string) (bool, error) {
	var exists bool
	if err := s.db.QueryRowContext(ctx, s.statements.PublicKeyExists, publicKey).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Store) MarkUsed(ctx context.Context, id int64, at time.Time, ip *string, userAgent *string) error {
	_, err := s.db.ExecContext(ctx, s.statements.MarkUsed, at.UTC(), stringPtrValue(ip), stringPtrValue(userAgent), id)
	return err
}

func (s *Store) Revoke(ctx context.Context, id int64, at time.Time) error {
	_, err := s.db.ExecContext(ctx, s.statements.Revoke, at.UTC(), id)
	return err
}

func (s *Store) findByID(ctx context.Context, id int64) (apikey.ApiKeyRecord, error) {
	row, err := s.scanOne(ctx, s.statements.FindByID, id)
	if err != nil {
		return apikey.ApiKeyRecord{}, err
	}
	return RowToRecord(row)
}

func (s *Store) scanOne(ctx context.Context, query string, args ...any) (ApiKeyRow, error) {
	var row ApiKeyRow
	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&row.ID,
		&row.App,
		&row.Env,
		&row.PublicKey,
		&row.Hash,
		&row.Scopes,
		&row.Name,
		&row.CreatedBy,
		&row.CreatedAt,
		&row.ExpiresAt,
		&row.RevokedAt,
		&row.LastUsedAt,
		&row.LastUsedIP,
		&row.LastUsedUserAgent,
	)
	return row, err
}
