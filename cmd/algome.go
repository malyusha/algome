package cmd

import (
	"fmt"
	"os"

	"github.com/caarlos0/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/malyusha/algome/genreadme"
	"github.com/urfave/cli/v2"
)

const (
	configFilePath = "./algome.conf.json"
)

var (
	configPath string
	debug      bool
	config     *genreadme.Config
)

var boldStyle = lipgloss.NewStyle().Bold(true)

// The root command for `algome` is the alias for sub command `generate`.
var app = &cli.App{
	Name:           "algome",
	Usage:          "CLI to manage solved problems readme",
	DefaultCommand: "generate",
	Before: func(ctx *cli.Context) error {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("using debug logs")
		}

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
		&cli.BoolFlag{
			Name:        "debug",
			Required:    false,
			Value:       false,
			Destination: &debug,
		},
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

func CreateDefaultConfig() *genreadme.Config {
	cfg := genreadme.NewConfig()
	cfg.Structure = genreadme.StructProviderConfig{
		Catalog: &genreadme.CatalogStructProviderConfig{
			BaseDir: "./",
		},
	}
	cfg.ProblemSources = []string{genreadme.ProviderLeetcode}

	return cfg
}
