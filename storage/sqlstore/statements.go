package sqlstore

type Statements struct {
	InsertAPIKey         string
	FindByID             string
	FindByPublicKey      string
	PublicKeyExists      string
	MarkUsed             string
	Revoke               string
}

func StatementsFor(dialect Dialect) Statements {
	return Statements{
		InsertAPIKey: `INSERT INTO api_keys (
    app, env, public_key, hash, scopes, name, created_by, created_at, expires_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`,
		FindByID: `SELECT
    id, app, env, public_key, hash, scopes, name, created_by, created_at,
    expires_at, revoked_at, last_used_at, last_used_ip, last_used_user_agent
FROM api_keys
WHERE id = ?
LIMIT 1;`,
		FindByPublicKey: `SELECT
    id, app, env, public_key, hash, scopes, name, created_by, created_at,
    expires_at, revoked_at, last_used_at, last_used_ip, last_used_user_agent
FROM api_keys
WHERE public_key = ?
LIMIT 1;`,
		PublicKeyExists: `SELECT EXISTS (SELECT 1 FROM api_keys WHERE public_key = ?);`,
		MarkUsed: `UPDATE api_keys
SET last_used_at = ?, last_used_ip = ?, last_used_user_agent = ?
WHERE id = ?;`,
		Revoke: `UPDATE api_keys
SET revoked_at = ?
WHERE id = ?
  AND revoked_at IS NULL;`,
	}
}
