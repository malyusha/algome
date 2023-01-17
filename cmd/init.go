package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/log"
	"github.com/malyusha/algome/logger"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:        "init",
	Description: "Initializes project",
	Action: func(ctx *cli.Context) error {
		filename := ctx.String("config")
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			logger.WithField("file", filename).Info(boldStyle.Render("creating config file"))
			return createConfigFile()
		}

		if err == nil {
			logger.Info(boldStyle.Render("configuration already initialized"))
			return nil
		}

		return fmt.Errorf("failed to check config file: %w", err)
	},
}

func createConfigFile() error {
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config to JSON: %w", err)
	}

	log.Info("configuration file successfully created")
	return nil
}
