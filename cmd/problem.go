package cmd

import (
	"github.com/caarlos0/log"
	"github.com/urfave/cli/v2"
)

var ProblemCommand = &cli.Command{
	Name:  "problem",
	Usage: "manage problems",
	Subcommands: cli.Commands{
		ProblemAddCommand,
	},
}

var ProblemAddCommand = &cli.Command{
	Name:  "add",
	Usage: "add new problem to solved list",
	Action: func(ctx *cli.Context) error {
		// todo:
		log.Info(boldStyle.Render("adding new problem"))
		return nil
	},
}
