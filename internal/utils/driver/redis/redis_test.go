package redis_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/outdead/goservice/internal/utils/driver/redis"
)

func TestNewClient(t *testing.T) {
	type TestCase struct {
		name    string
		config  redis.Config
		wantErr bool
	}

	tests := []TestCase{
		{"empty addr", redis.Config{}, true},
		{"wrong addr", redis.Config{Addr: "localhost:4444"}, true},
	}

	if run := getVar("TEST_REAL_REDIS", "false"); run == "true" {
		tests = append(tests, TestCase{
			name: "real db positive",
			config: redis.Config{
				Addr: getVar("TEST_REDIS_ADDR", "127.0.0.1:6379"),
				DB:   getIntVar("TEST_REDIS_DB", 0),
			},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := redis.NewClient(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("cteate client error expected: %v, got %v", tt.wantErr, err)
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

func getIntVar(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		res, err := strconv.Atoi(value)
		if err == nil {
			// We don't care about mistakes here. Didn't convert - returned the
			// transmitted default.
			return res
		}
	}

	return fallback
}
