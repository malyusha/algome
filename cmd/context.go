package cmd

import (
	"github.com/malyusha/algome/logger"
)

type Context struct {
	logger         logger.Logger
	configFilepath string
}
