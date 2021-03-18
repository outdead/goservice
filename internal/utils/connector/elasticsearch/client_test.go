package elasticsearch_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/outdead/goservice/internal/utils/connector/elasticsearch"
)

// FakeModel implements Model interface.
type FakeModel struct {
	ID   int64  `json:"id"`
	Data string `json:"data"`
}

func (m *FakeModel) TableName() string {
	return "test"
}

func (m *FakeModel) CalculateID() string {
	return fmt.Sprintf("%d", m.ID)
}

func TestClient_MultiInsert(t *testing.T) {
	if run := getVar("TEST_REAL_ELASTIC", "false"); run == "true" {
		t.Run("real db elastic", func(t *testing.T) {
			cfg := elasticsearch.Config{
				Addr:     getVar("TEST_ELASTIC_ADDR", "http://localhost:9200"),
				Database: "connector_test",
			}

			client, err := elasticsearch.NewClient(&cfg)
			if err != nil {
				t.Fatal(err)
			}

			defer client.Close()

			// Cleanup.
			defer client.Conn().DeleteIndex(cfg.Database).Do(context.Background())

			wantPost := FakeModel{ID: 1231, Data: "data 1"}
			wantID := "1231"

			posts := []elasticsearch.Model{
				&wantPost,
				&FakeModel{ID: 1232, Data: "data 2"},
				&FakeModel{ID: 1233, Data: "data 3"},
			}

			if err := client.MultiInsert(posts); err != nil {
				t.Fatal(err)
			}

			// Get one of the inserted records.
			typ := new(FakeModel).TableName()
			get, err := client.Conn().Get().Index(cfg.Database).Type(typ).Id(wantID).Do(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			// If you have not received a record by the generated identifier, exit.
			if !get.Found {
				t.Errorf("record with id was not found: %q", wantID)
			}

			// Now let's compare the recorded and received value.
			postBytes, err := get.Source.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}

			wantPostBytes, _ := json.Marshal(wantPost)
			if string(wantPostBytes) != string(postBytes) {
				t.Errorf("\ngot:  %s \nwant: %s", postBytes, wantPostBytes)
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
