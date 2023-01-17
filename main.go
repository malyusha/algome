package main

import (
	"os"

	"github.com/malyusha/algome/cmd"
)

func main() {
	var command string
	args := os.Args[1:]
	if len(args) != 0 {
		command = args[0]
	}
	cmd.Execute(command)
}
