CREATE TABLE IF NOT EXISTS api_keys (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    app VARCHAR(3) NOT NULL,
    env VARCHAR(10) NOT NULL,
    public_key CHAR(16) NOT NULL,
    hash CHAR(64) NOT NULL,
    scopes TEXT NOT NULL,
    name VARCHAR(120) NULL,
    created_by VARCHAR(180) NULL,
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    expires_at TIMESTAMP(6) NULL,
    revoked_at TIMESTAMP(6) NULL,
    last_used_at TIMESTAMP(6) NULL,
    last_used_ip VARCHAR(45) NULL,
    last_used_user_agent VARCHAR(512) NULL,
    CONSTRAINT ux_api_keys_public_key UNIQUE (public_key),
    INDEX ix_api_keys_app_env (app, env),
    INDEX ix_api_keys_revoked_at (revoked_at),
    INDEX ix_api_keys_expires_at (expires_at)
);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) NOT NULL PRIMARY KEY,
    applied_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);
