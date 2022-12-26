package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/malyusha/algome/logger"
)

const (
	configFilePath = "./algome.conf.json"
)

func Execute(cmd string) {
	var ctx Context
	flag.StringVar(&ctx.configFilepath, "conf", configFilePath, "Configuration file path. Default - $CUR_DIR/algome.conf.json")

	flag.Parse()
	log := logger.NewDefaultSimpleLogger()
	logger.SetGlobalLogger(log)

	ctx.logger = log

	var err error
	switch cmd {
	case "", "generate":
		err = generateCommand(ctx)
	case "init":
		err = initCommand(ctx)
	default:
		err = fmt.Errorf("unknown command '%s'", cmd)
		flag.Usage()
	}

	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(1)
	}
}
