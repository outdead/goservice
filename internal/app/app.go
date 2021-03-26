package app

import (
	"fmt"
	"os"

	"github.com/outdead/goservice/internal/app/daemon"
	"github.com/outdead/goservice/internal/utils/logutil"
	"github.com/urfave/cli/v2"
)

// App is main application.
type App struct {
	name    string
	version string
	logger  *logutil.Logger
	cli     *cli.App
}

// New creates and returns new App.
func New(name, version string) *App {
	app := App{
		name:    name,
		version: version,
		logger:  logutil.New(logutil.SetService(name), logutil.SetVersion(version)),
	}

	return &app
}

// Run executes main app process.
func (a *App) Run() {
	a.init()

	if err := a.cli.Run(os.Args); err != nil {
		a.logger.NewEntry().Fatal(err)
	}
}

func (a *App) init() {
	app := cli.NewApp()
	app.Name = a.name
	app.Version = a.version
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

	app.Action = a.action()

	a.cli = app
}

func (a *App) action() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cfg, err := daemon.NewConfig(c.String("config"))
		if err != nil {
			return fmt.Errorf("new config: %w", err)
		}

		if c.Bool("print") {
			a.logger.NewEntry().Info("got -p flag - print config and terminate")

			return cfg.Print()
		}

		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validate config: %w", err)
		}

		a.logger.Customize(&cfg.App.Log)

		d := daemon.NewDaemon(cfg, a.logger.NewEntry())

		defer func() {
			if err := d.Close(); err != nil {
				a.logger.NewEntry().Errorf("close daemon err: %s", err)
			}
		}()

		return d.Run()
	}
}
