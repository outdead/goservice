package postgres

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	type TestCase struct {
		name    string
		config  Config
		wantErr bool
	}

	tests := []TestCase{
		{"empty password", Config{Addr: "localhost:5432", Database: "goservice", User: "postgres"}, true},
	}

	if run := getVar("TEST_REAL_POSTGRES", "false"); run == "true" {
		tests = append(tests, TestCase{
			name: "real db positive",
			config: Config{
				Addr:     getVar("TEST_POSTGRES_ADDR", "127.0.0.1:5432"),
				Database: getVar("TEST_POSTGRES_DB", "goservice"),
				User:     getVar("TEST_POSTGRES_USER", "postgres"),
				Password: getVar("TEST_POSTGRES_PASSWORD", "postgres"),
				PoolSize: 10,
				Debug:    false,
			},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDB(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("cteate db error expected: %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func getVar(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
