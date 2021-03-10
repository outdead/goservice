package main

import (
	"fmt"
	"os"

	"github.com/outdead/goservice/internal/daemon"
	"github.com/outdead/goservice/internal/utils/logutil"
	"github.com/urfave/cli/v2"
)

// ServiceName contains the name of the service. Displayed in logs and when help
// command is called.
const ServiceName = "goservice"

// ServiceVersion contains the service version number in the semantic versioning
// format (http://semver.org/). Displayed in logs and when help command is
// called. During service compilation, you can pass the version value using the
// flag `-ldflags "-X main.Version=${VERSION}"`.
var ServiceVersion = "0.0.0-develop"

func main() {
	log := logutil.New(logutil.SetService(ServiceName), logutil.SetVersion(ServiceVersion))

	app := cli.NewApp()
	app.Name = ServiceName
	app.Version = ServiceVersion
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
			Usage:    "Path to config file",
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "Print config file and exit",
		},
	}

	app.Action = action(log)

	if err := app.Run(os.Args); err != nil {
		log.NewEntry().Fatal(err)
	}
}

func action(log *logutil.Logger) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cfg, err := daemon.NewConfig(c.String("config"))
		if err != nil {
			return fmt.Errorf("new config: %w", err)
		}

		if c.Bool("print") {
			log.NewEntry().Info("got -p flag - print config and terminate")

			return cfg.Print()
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate config: %w", err)
		}

		log.Customize(&cfg.App.Log)

		d := daemon.NewDaemon(cfg, log.NewEntry())

		defer func() {
			if err := d.Close(); err != nil {
				log.NewEntry().Errorf("close daemon err: %s", err)
			}
		}()

		return d.Run()
	}
}
