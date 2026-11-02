package platform

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Simulate env vars
	os.Setenv("DB_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("PORT", "1234")
	os.Setenv("JWT_SECRET", "secret")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "1234" {
		t.Errorf("expected port 1234, got %s", cfg.Port)
	}
}
