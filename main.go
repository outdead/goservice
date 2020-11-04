package main

import (
	"os"

	"github.com/outdead/echo-skeleton/internal/daemon"
	"github.com/outdead/echo-skeleton/internal/logger"
	"github.com/urfave/cli/v2"
)

// ServiceName содержит имя сервиса. Выводится в логах и при вызове help.
const ServiceName = "golang echo skeleton"

// ServiceVersion содержит номер версии сервиса в формате семантическо
// версионирования (http://semver.org/). Выводится в логах и при вызове help.
// Во время компиляции сервиса можно передать значение версии при помощи
// флага `-ldflags "-X main.Version=${VERSION}"`.
var ServiceVersion = "0.0.0-develop"

func main() {
	log := logger.New(logger.SetService(ServiceName), logger.SetVersion(ServiceVersion))

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
		log.WithAppInfo().Fatal(err)
	}
}

func action(log *logger.Logger) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cfg, err := daemon.NewConfig(c.String("config"))
		if err != nil {
			return err
		}

		if c.Bool("print") {
			log.WithAppInfo().Info("got -p flag - print config and terminate")
			return cfg.Print()
		}

		if err := cfg.Validate(); err != nil {
			return err
		}

		log.Customize(&cfg.App.Log)

		d := daemon.NewDaemon(cfg, log.WithAppInfo())

		defer func() {
			if err := d.Close(); err != nil {
				log.WithAppInfo().Errorf("close daemon err: %s", err)
			}
		}()

		return d.Run()
	}
}
