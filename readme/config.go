package readme

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/malyusha/algome/logger"
	"github.com/malyusha/algome/readme/generator"
)

type Config struct {
	OutputDirectory      string               `json:"output"`
	SolutionsBaseDir     string               `json:"solutions_dir"`
	Structure            StructProviderConfig `json:"structure"`
	ProblemSources       []string             `json:"sources"`
	ProvidersCacheDir    string               `json:"problems_cache_dir"`
	OverloadTemplatesDir *string              `json:"templates_dir,omitempty"`

	templates *generator.Templates
}

var DefaultConfig = CreateDefaultConfig()

func CreateDefaultConfig() *Config {
	cfg := &Config{
		OutputDirectory:  "./",
		SolutionsBaseDir: "./",
		Structure: StructProviderConfig{
			Catalog: &CatalogStructProviderConfig{
				BaseDir: "./",
			},
		},
		ProvidersCacheDir:    ".cache",
		ProblemSources:       []string{providerLeetcode},
		OverloadTemplatesDir: nil,
		templates:            &generator.Templates{},
	}

	usr, err := user.Current()
	if err == nil {
		cfg.ProvidersCacheDir = filepath.Join(usr.HomeDir, ".algome/cache")
	}

	return cfg
}

type FileReadmeConfig struct {
	OutputFileBasename string `json:"output_file_basename"`
}

type CatalogStructProviderConfig struct {
	MapAttribute string `json:"map_attr"`
	BaseDir      string `json:"base_dir"`
}

type StructProviderConfig struct {
	Catalog *CatalogStructProviderConfig `json:"catalog"`
}

func Load(configFile string) (*Config, error) {
	cfg := *DefaultConfig

	if _, err := os.Stat(configFile); err == nil && !os.IsNotExist(err) {
		b, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("read config file error: %w", err)
		}

		if err := json.Unmarshal(b, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal config error: %w", err)
		}
	} else {
		logger.Warn("can't load configuration file. using default configuration")
	}

	err := generator.LoadTemplates(cfg.templates, generator.DefaultTemplates)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates")
	}

	if cfg.OverloadTemplatesDir != nil {
		if err := overloadTemplates(*cfg.OverloadTemplatesDir, cfg.templates); err != nil {
			logger.Error("failed to overload templates: %s", err.Error())
		}
	}

	return &cfg, nil
}

func overloadTemplates(dir string, base *generator.Templates) error {
	overloadList := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			// skip directories
			return nil
		}

		filename := filepath.Base(info.Name())
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)

		fileContents, err := os.ReadFile(info.Name())
		if err != nil {
			return fmt.Errorf("ReadFile error: %w", err)
		}

		overloadList[name] = string(fileContents)
		return nil
	})

	if err != nil {
		return fmt.Errorf("overloadTemplates walk error: %w", err)
	}

	return generator.LoadTemplates(base, overloadList)
}
