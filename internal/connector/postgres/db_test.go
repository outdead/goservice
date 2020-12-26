package postgres

import (
	"testing"
)

func TestNewDB(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		// {"positive", config,false},
		{"empty password", Config{Addr: "localhost:5432", Database: "pcs", User: "postgres"}, true},
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
