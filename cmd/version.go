package cmd

import (
	"fmt"

	"github.com/malyusha/algome/readme/version"
	"github.com/urfave/cli/v2"
)

var VersionCommand = &cli.Command{
	Name:  "version",
	Usage: "Print current version of CLI",
	Action: func(ctx *cli.Context) error {
		fmt.Println(version.String())
		return nil
	},
}
