package migrations_test

import (
	"os"
	"strings"
	"testing"
)

func TestMigration001ContainsExpectedSchema(t *testing.T) {
	data, err := os.ReadFile("001_create_webhook_events.sql")
	if err != nil {
		t.Fatalf("TestMigration001ContainsExpectedSchema: %v", err)
	}

	sql := string(data)

	required := []struct {
		token string
		desc  string
	}{
		{"CREATE TABLE webhook_events", "table creation"},
		{"delivery_id", "delivery_id column"},
		{"event_type", "event_type column"},
		{"payload", "payload column"},
		{"status", "status column"},
		{"created_at", "created_at column"},
		{"processed_at", "processed_at column"},
		{"UNIQUE", "unique constraint on delivery_id"},
		{"CREATE INDEX", "status+created_at index for cron query"},
	}

	for _, r := range required {
		if !strings.Contains(sql, r.token) {
			t.Errorf("TestMigration001ContainsExpectedSchema: migration missing %s (expected %q)", r.desc, r.token)
		}
	}
}
