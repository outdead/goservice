package elasticsearch_test

import (
	"testing"

	"github.com/outdead/goservice/internal/utils/connector/elasticsearch"
)

var config = elasticsearch.Config{
	Addr:                "http://localhost:9200",
	Database:            "test",
	HealthcheckInterval: elasticsearch.DefaultBatchLimit,
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  elasticsearch.Config
		wantErr bool
	}{
		{"positive validation", config, false},
		{"empty addr", elasticsearch.Config{}, true},
		{"empty database", elasticsearch.Config{Addr: "http://localhost:9200"}, true},
		{"negative healthcheck_interval", elasticsearch.Config{Addr: "http://localhost:9200", HealthcheckInterval: -1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("validation error expected: %v, got %v", tt.wantErr, err)
			}
		})
	}
}
