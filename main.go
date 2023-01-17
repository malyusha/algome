package main

import (
	"os"

	"github.com/malyusha/algome/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
