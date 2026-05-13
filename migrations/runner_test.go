package migrations

import "testing"

func TestRunnerDefaults(t *testing.T) {
	r := NewRunner(nil, "")
	if r.migrationsDir == "" {
		t.Fatal("expected default migrations dir")
	}
}
