package tickerprocess_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/outdead/goservice/internal/utils/logutil"
	"github.com/outdead/goservice/snippets/tickerprocess"
)

type FakeRepo struct {
	internal struct {
		// Add here interface functions.
	}
}

func NewFakeRepo() *FakeRepo {
	repo := FakeRepo{}

	return &repo
}

func TestNewProcess(t *testing.T) {
	output := &bytes.Buffer{}

	cfg := tickerprocess.Config{
		StartInterval: time.Duration(1 * time.Second),
	}

	log := logutil.New()
	log.SetLevel("debug")
	log.SetOutput(output)

	// Inject interface functions.
	repo := NewFakeRepo()

	processor := tickerprocess.NewProcess(&cfg, repo, log.NewEntry())

	processor.Run()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-ticker.C:
			log.Info("time left")

			break Loop
		case err := <-processor.Errors():
			log.WithError(err).Error("got error from process")

			break Loop
		}
	}

	processor.Quit()

	logs := string(output.Bytes())

	// Look for processing.
	want := "tick..."
	if !strings.Contains(logs, want) {
		t.Errorf("logs message was not found: %q", want)

		fmt.Println(string(output.Bytes()))
	}

	// Make sure that the termination has occurred correctly
	want = "process stopped"
	if !strings.Contains(logs, want) {
		t.Errorf("logs message was not found: %q", want)

		fmt.Println(string(output.Bytes()))
	}
}
