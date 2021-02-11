package process_test

import (
	"time"

	"github.com/outdead/goservice/internal/utils/logutils"
	"github.com/outdead/goservice/snippets/process"
)

func ExampleNewProcess() {
	cfg := process.Config{
		StartInterval: time.Duration(1 * time.Second),
	}

	log := logutils.New()
	log.SetLevel("debug")

	processor := process.NewProcess(&cfg, log.NewEntry(), nil)

	go processor.Run()

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

	// Output:
}
