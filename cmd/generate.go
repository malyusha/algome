package cmd

import (
	"context"
	"fmt"

	"github.com/caarlos0/log"
	"github.com/malyusha/algome/readme"
	"github.com/urfave/cli/v2"
)

var GenerateCommand = &cli.Command{
	Name:        "generate",
	Description: "Generates new readme",
	Action: func(ctx *cli.Context) error {
		gen, err := readme.NewGenerator(config)
		if err != nil {
			return fmt.Errorf("failed to create generator: %w", err)
		}
		if err := gen.GenerateReadme(context.Background()); err != nil {
			return err
		}

		log.Info("readme successfully generated")

		return nil
	},
}
