package postgres

import (
	"testing"
)

var config = Config{
	Addr:     "127.0.0.1:5432",
	Database: "goservice",
	User:     "postgres",
	Password: "postgres",
	PoolSize: 10,
	Debug:    true,
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{"positive validation", config, false},
		{"empty addr", Config{}, true},
		{"empty database", Config{Addr: "localhost:5432"}, true},
		{"empty user", Config{Addr: "localhost:5432", Database: "goservice"}, true},
		{"empty password", Config{Addr: "localhost:5432", Database: "goservice", User: "postgres"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validation error expected: %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestConfig_GetDataSourceNameN(t *testing.T) {
	expected := "postgres://postgres:postgres@127.0.0.1:5432/goservice?sslmode=disable"

	if got := config.GetDataSourceName(); got != expected {
		t.Errorf("dns expected: %q, got %q", expected, got)
	}
}
