package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:        "init",
	Description: "Initializes project",
	Action: func(ctx *cli.Context) error {
		_, err := os.Stat(ctx.String("config"))
		if os.IsNotExist(err) {
			return createConfigFile(ctx)
		}

		if err == nil {
			return nil
		}

		return fmt.Errorf("failed to check config file: %w", err)
	},
}

func createConfigFile(ctx *cli.Context) error {
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config to JSON: %w", err)
	}

	getLogger(ctx).Info("configuration file successfully created")
	return nil
}
