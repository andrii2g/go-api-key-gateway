package apikey

import "errors"

var (
	ErrInvalidApp            = errors.New("invalid app prefix")
	ErrInvalidEnv            = errors.New("invalid environment")
	ErrInvalidScope          = errors.New("invalid scope")
	ErrInvalidSecretBytes    = errors.New("invalid secret byte count")
	ErrInvalidPepper         = errors.New("invalid pepper")
	ErrInvalidUsageQueueSize = errors.New("invalid usage queue size")
	ErrInvalidTimeout        = errors.New("invalid timeout")
	ErrPublicKeyCollision    = errors.New("public key collision retry limit exceeded")
)

type ValidationFailureReason string

const (
	FailureNone                ValidationFailureReason = "None"
	FailureMissing             ValidationFailureReason = "Missing"
	FailureMalformed           ValidationFailureReason = "Malformed"
	FailureInvalid             ValidationFailureReason = "Invalid"
	FailureEnvironmentMismatch ValidationFailureReason = "EnvironmentMismatch"
	FailureExpired             ValidationFailureReason = "Expired"
	FailureRevoked             ValidationFailureReason = "Revoked"
	FailureScopeDenied         ValidationFailureReason = "ScopeDenied"
	FailureStoreUnavailable    ValidationFailureReason = "StoreUnavailable"
)

type ValidationResult struct {
	OK        bool
	Principal *Principal
	Reason    ValidationFailureReason
}
