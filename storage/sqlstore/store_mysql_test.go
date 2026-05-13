package sqlstore

import (
	"context"
	"os"
	"testing"
)

func TestMySQLIntegrationConfigured(t *testing.T) {
	if os.Getenv("TEST_MYSQL_DSN") == "" {
		t.Skip("TEST_MYSQL_DSN is not set")
	}
	ctx := context.Background()
	db, err := Open(ctx, Options{DSN: os.Getenv("TEST_MYSQL_DSN")})
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	_ = db.Close()
}
