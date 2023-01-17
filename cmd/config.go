package cmd

import (
	"github.com/malyusha/algome/readme"
)

func CreateDefaultConfig() *readme.Config {
	cfg := readme.NewConfig()
	cfg.Structure = readme.StructProviderConfig{
		Catalog: &readme.CatalogStructProviderConfig{
			BaseDir: "./",
		},
	}
	cfg.ProblemSources = []string{readme.ProviderLeetcode}

	return cfg
}
