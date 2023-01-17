package readme

import (
	"errors"
	"fmt"

	"github.com/malyusha/algome/logger"
	"github.com/malyusha/algome/readme/generator"
	"github.com/malyusha/algome/readme/sourceproviders"
	"github.com/malyusha/algome/readme/structproviders"
)

const (
	ProviderLeetcode = "leetcode"
)

func NewGenerator(cfg *Config) (*generator.Generator, error) {
	return newGeneratorFromConfig(cfg)
}

func newGeneratorFromConfig(cfg *Config) (*generator.Generator, error) {
	structProvider, err := structProviderFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create struct provider error: %w", err)
	}

	providers := sourceProvidersFromConfig(cfg)

	gen := generator.NewGenerator(cfg.outputDir(), cfg.templates, structProvider, providers)

	return gen, nil
}

func sourceProvidersFromConfig(cfg *Config) []generator.Provider {
	out := make([]generator.Provider, 0)
	for _, source := range cfg.ProblemSources {
		sp, err := createSourceProvider(source)
		if err != nil {
			logger.Error("failed to create source provider: %s", err.Error())
			continue
		}

		namedProvider := generator.NewProvider(source, cfg.providersCacheDir(), sp)
		out = append(out, namedProvider)
	}

	return out
}

func createSourceProvider(source string) (generator.ProblemsProvider, error) {
	switch source {
	case ProviderLeetcode:
		return sourceproviders.NewLeetcodeProvider()
	}

	return nil, fmt.Errorf("unsupported provider '%s'", source)
}

func structProviderFromConfig(cfg *Config) (generator.StructureProvider, error) {
	var (
		provider generator.StructureProvider
		err      error
	)
	switch {
	case cfg.Structure.Catalog != nil:
		provider = structproviders.NewCatalogStructProvider(
			structproviders.ProblemIdType(cfg.Structure.Catalog.MapAttribute),
			cfg.Structure.Catalog.BaseDir,
		)
	}

	if provider == nil {
		return nil, errors.New("no known struct provider specified")
	}

	return provider, err
}
