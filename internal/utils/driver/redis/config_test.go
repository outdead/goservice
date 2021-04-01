package redis_test

import (
	"testing"
	"time"

	"github.com/outdead/goservice/internal/utils/driver/redis"
)

var config = redis.Config{
	Addr:         "127.0.0.1:6379",
	Password:     "",
	DB:           1,
	TTL:          5 * time.Second,
	MaxRetries:   1,
	DialTimeout:  0,
	ReadTimeout:  0,
	WriteTimeout: 0,
	PoolSize:     10,
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  redis.Config
		wantErr bool
	}{
		{"positive validation", config, false},
		{"empty addr", redis.Config{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validation error expected: %v, got %v", tt.wantErr, err)
			}
		})
	}
}
