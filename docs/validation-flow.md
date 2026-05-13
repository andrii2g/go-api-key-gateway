# Validation Flow

1. Parse the raw bearer token.
2. Lookup by `public_key`.
3. Reject environment mismatch, revocation, expiry, and hash mismatch.
4. Check required scopes.
5. Return a principal and enqueue best-effort usage metadata.
