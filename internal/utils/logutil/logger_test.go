package logutil

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Message struct {
	Level   string `json:"level,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Service string `json:"service,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

func TestNewLogger(t *testing.T) {
	expected := `{"level":"warning","msg":"hello","service":"test-service","version":"0.0.0"}`
	output := &bytes.Buffer{}

	opts := []Option{
		SetLevel("debug"),
		SetVersion("0.0.0"),
		SetService("test-service"),
	}

	l := New(opts...)
	l.SetOutput(output)

	l.NewEntry().Warning("hello")

	var msg Message
	if err := json.Unmarshal(output.Bytes(), &msg); err != nil {
		t.Fatal(err)
	}

	js, _ := json.Marshal(msg)
	got := string(js)

	assert.Equal(t, expected, got, "wrong fields")

	l.Customize(&Config{Level: "notexists"})
	assert.Equal(t, l.Level.String(), "info")
}
