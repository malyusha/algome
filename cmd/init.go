package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/malyusha/algome/readme"
)

func initCommand(ctx Context) error {
	_, err := os.Stat(ctx.configFilepath)
	if os.IsNotExist(err) {
		return createConfigFile(ctx)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to check config file: %w", err)
}

func createConfigFile(ctx Context) error {
	file, err := os.Create(ctx.configFilepath)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %w", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(readme.DefaultConfig); err != nil {
		return fmt.Errorf("failed to encode config to JSON: %w", err)
	}

	ctx.logger.Info("configuration file successfully created")
	return nil
}
