package apikey

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"
)

type Store interface {
	Create(ctx context.Context, record ApiKeyRecord) (ApiKeyRecord, error)
	FindByPublicKey(ctx context.Context, publicKey string) (*ApiKeyRecord, error)
	PublicKeyExists(ctx context.Context, publicKey string) (bool, error)
	MarkUsed(ctx context.Context, id int64, at time.Time, ip *string, userAgent *string) error
	Revoke(ctx context.Context, id int64, at time.Time) error
}

type Service struct {
	store   Store
	options Options
	usage   UsageRecorder
	now     func() time.Time
}

func NewService(store Store, options Options, usage UsageRecorder) (*Service, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	if usage == nil {
		usage = NoopUsageRecorder{}
	}
	return &Service{
		store:   store,
		options: options.Normalize(),
		usage:   usage,
		now:     time.Now,
	}, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (CreateResult, error) {
	log.Printf("apikey create start app=%q env=%q", req.App, req.Env)
	app, err := NormalizeApp(req.App)
	if err != nil {
		return CreateResult{}, err
	}
	log.Printf("apikey create normalized app=%q", app)
	env, err := NormalizeEnv(req.Env)
	if err != nil {
		return CreateResult{}, err
	}
	log.Printf("apikey create normalized env=%q", env)
	scopes, err := NormalizeScopes(req.Scopes)
	if err != nil {
		return CreateResult{}, err
	}
	log.Printf("apikey create normalized scopes=%v", scopes)
	if req.ExpiresAt != nil && !req.ExpiresAt.After(s.now().UTC()) {
		return CreateResult{}, ErrInvalidEnv
	}
	log.Printf("apikey create expires_at validated")

	var publicKey string
	for attempt := 0; attempt < 10; attempt++ {
		log.Printf("apikey create generating public key attempt=%d", attempt+1)
		candidate, err := GeneratePublicKey()
		if err != nil {
			return CreateResult{}, err
		}
		log.Printf("apikey create checking public key collision attempt=%d public_key=%s", attempt+1, candidate)
		exists, err := s.store.PublicKeyExists(ctx, candidate)
		if err != nil {
			return CreateResult{}, err
		}
		if !exists {
			publicKey = candidate
			break
		}
	}
	if publicKey == "" {
		return CreateResult{}, ErrPublicKeyCollision
	}

	secret, err := GenerateSecret(s.options.SecretBytes)
	if err != nil {
		return CreateResult{}, err
	}
	log.Printf("apikey create generated secret and hash for public_key=%s", publicKey)
	hash, err := HashSecret(secret, s.options.Pepper)
	if err != nil {
		return CreateResult{}, err
	}

	now := s.now().UTC()
	record := ApiKeyRecord{
		App:       app,
		Env:       env,
		PublicKey: publicKey,
		Hash:      hash,
		Scopes:    scopes,
		Name:      trimOptionalString(req.Name),
		CreatedBy: trimOptionalString(req.CreatedBy),
		CreatedAt: now,
		ExpiresAt: utcTimePtr(req.ExpiresAt),
	}
	saved, err := s.store.Create(ctx, record)
	if err != nil {
		return CreateResult{}, err
	}
	log.Printf("apikey create stored id=%d public_key=%s", saved.ID, saved.PublicKey)

	return CreateResult{
		ID:        saved.ID,
		App:       saved.App,
		Env:       saved.Env,
		PublicKey: saved.PublicKey,
		APIKey:    BuildFullKey(saved.App, saved.PublicKey, secret),
		Scopes:    saved.Scopes,
		Name:      saved.Name,
		CreatedBy: saved.CreatedBy,
		CreatedAt: saved.CreatedAt,
		ExpiresAt: saved.ExpiresAt,
	}, nil
}

func (s *Service) Validate(ctx context.Context, req ValidationRequest) ValidationResult {
	parsed, reason := Parse(req.RawKey)
	if reason == FailureMissing {
		return ValidationResult{Reason: FailureMissing}
	}
	if reason != FailureNone {
		return ValidationResult{Reason: FailureMalformed}
	}

	env, err := NormalizeEnv(req.Environment)
	if err != nil {
		return ValidationResult{Reason: FailureEnvironmentMismatch}
	}
	record, err := s.store.FindByPublicKey(ctx, parsed.PublicKey)
	if err != nil {
		return ValidationResult{Reason: FailureStoreUnavailable}
	}
	if record == nil {
		return ValidationResult{Reason: FailureInvalid}
	}
	if record.App != parsed.App {
		return ValidationResult{Reason: FailureInvalid}
	}
	if record.Env != env {
		return ValidationResult{Reason: FailureEnvironmentMismatch}
	}
	if record.RevokedAt != nil {
		return ValidationResult{Reason: FailureRevoked}
	}
	now := s.now().UTC()
	if record.ExpiresAt != nil && !record.ExpiresAt.After(now) {
		return ValidationResult{Reason: FailureExpired}
	}

	ok, err := CompareSecretHash(parsed.Secret, record.Hash, s.options.Pepper)
	if err != nil || !ok {
		return ValidationResult{Reason: FailureInvalid}
	}
	if !HasRequiredScopes(record.Scopes, req.RequiredScopes) {
		return ValidationResult{Reason: FailureScopeDenied}
	}

	principal := &Principal{
		KeyID:     record.ID,
		App:       record.App,
		Env:       record.Env,
		PublicKey: record.PublicKey,
		Scopes:    append([]string(nil), record.Scopes...),
		Name:      record.Name,
	}
	if s.usage != nil {
		s.usage.Record(UsageEvent{
			KeyID:     record.ID,
			At:        now,
			IP:        req.IP,
			UserAgent: req.UserAgent,
		})
	}
	return ValidationResult{
		OK:        true,
		Principal: principal,
		Reason:    FailureNone,
	}
}

func (s *Service) Revoke(ctx context.Context, id int64, at time.Time) error {
	return s.store.Revoke(ctx, id, at.UTC())
}

var (
	appPattern = regexp.MustCompile(`^[a-z0-9]{1,3}$`)
	envPattern = regexp.MustCompile(`^[a-z0-9-]{1,10}$`)
)

func NormalizeApp(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if !appPattern.MatchString(normalized) {
		return "", ErrInvalidApp
	}
	return normalized, nil
}

func NormalizeEnv(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if !envPattern.MatchString(normalized) {
		return "", ErrInvalidEnv
	}
	return normalized, nil
}

func isValidApp(value string) bool {
	return appPattern.MatchString(value)
}

func isValidPublicKey(value string) bool {
	if len(value) != PublicKeyLength {
		return false
	}
	for _, r := range value {
		if !strings.ContainsRune(PublicKeyAlphabet, r) {
			return false
		}
	}
	return true
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func utcTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	t := value.UTC()
	return &t
}
