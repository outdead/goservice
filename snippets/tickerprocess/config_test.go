package tickerprocess_test

import (
	"testing"
	"time"

	"github.com/outdead/goservice/snippets/tickerprocess"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  tickerprocess.Config
		wantErr bool
	}{
		{"positive validation", tickerprocess.Config{
			StartInterval: 10 * time.Second,
		}, false},
		{"disabled config", tickerprocess.Config{
			Disabled: true,
		}, false},
		{"empty config", tickerprocess.Config{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validation error expected: %v, got %v", tt.wantErr, err)
			}
		})
	}
}
