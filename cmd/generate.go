package cmd

import (
	"context"
	"fmt"

	"github.com/malyusha/algome/readme"
)

func generateCommand(ctx Context) error {
	gen, err := readme.NewGenerator(ctx.configFilepath)
	if err != nil {
		return fmt.Errorf("failed to create generator: %w", err)
	}
	if err := gen.GenerateReadme(context.Background()); err != nil {
		return err
	}

	ctx.logger.Info("readme successfully generated")

	return nil
}
