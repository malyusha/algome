package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/malyusha/algome/logger"
	"github.com/malyusha/algome/readme"
	"github.com/urfave/cli/v2"
)

const (
	configFilePath = "./algome.conf.json"
)

var (
	configPath string
	config     *readme.Config
)

// The root command for `algome` is the alias for sub command `generate`.
var app = &cli.App{
	Name:           "algome",
	Usage:          "CLI to manage solved problems readme",
	DefaultCommand: "generate",
	Before: func(ctx *cli.Context) error {
		log := logger.NewDefaultSimpleLogger(logger.LevelInfo)
		logger.SetGlobalLogger(log)
		ctx.Context = context.WithValue(ctx.Context, "logger", log)

		config = CreateDefaultConfig()
		configPath := ctx.String("config")
		if configPath != "" {
			err := config.Load(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
		}
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Required:    false,
			Value:       configFilePath,
			TakesFile:   true,
			Destination: &configPath,
			Action: func(c *cli.Context, s string) error {
				fmt.Println("action config parse", s)
				return nil
			},
		},
	},
}

func init() {
	app.Commands = []*cli.Command{
		GenerateCommand,
		InitCommand,
		ProblemCommand,
	}
}

func Execute() error {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(app.ErrWriter, err.Error())
	}
	return nil
}

func getLogger(ctx *cli.Context) logger.Logger {
	return ctx.Context.Value("logger").(logger.Logger)
}

func CreateDefaultConfig() *readme.Config {
	cfg := readme.NewConfig()
	cfg.Structure = readme.StructProviderConfig{
		Catalog: &readme.CatalogStructProviderConfig{
			BaseDir: "./",
		},
	}
	cfg.ProblemSources = []string{readme.ProviderLeetcode}

	return cfg
}
