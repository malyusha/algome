package readme

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/malyusha/algome/logger"
	"github.com/malyusha/algome/readme/generator"
)

const (
	defaultProvidersCacheDir = "~/.algome/cache"
	defaultOutputDir         = "./"
)

type Config struct {
	OutputDirectory      string               `json:"output,omitempty"`
	Structure            StructProviderConfig `json:"structure"`
	ProblemSources       []string             `json:"sources"`
	ProvidersCacheDir    string               `json:"problems_cache_dir,omitempty"`
	OverloadTemplatesDir *string              `json:"templates_dir,omitempty"`

	templates *generator.Templates
}

func (c *Config) outputDir() string {
	if c.OutputDirectory != "" {
		return c.OutputDirectory
	}

	return defaultOutputDir
}

func (c *Config) providersCacheDir() string {
	if c.ProvidersCacheDir != "" {
		return c.ProvidersCacheDir
	}

	return defaultProvidersCacheDir
}

func NewConfig() *Config {
	return &Config{
		templates: &generator.Templates{},
	}
}

func (c *Config) SetTemplates(t *generator.Templates) {
	c.templates = &generator.Templates{}
}

type FileReadmeConfig struct {
	OutputFileBasename string `json:"output_file_basename"`
}

type CatalogStructProviderConfig struct {
	MapAttribute string `json:"map_attr,omitempty"`
	BaseDir      string `json:"base_dir"`
}

type StructProviderConfig struct {
	Catalog *CatalogStructProviderConfig `json:"catalog"`
}

func (c *Config) Load(configFile string) error {
	if _, err := os.Stat(configFile); err == nil && !os.IsNotExist(err) {
		b, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("read config file error: %w", err)
		}

		if err := json.Unmarshal(b, &c); err != nil {
			return fmt.Errorf("unmarshal config error: %w", err)
		}
	}

	err := generator.LoadTemplates(c.templates, generator.DefaultTemplates)
	if err != nil {
		return fmt.Errorf("failed to load templates")
	}

	if c.OverloadTemplatesDir != nil {
		if err := overloadTemplates(*c.OverloadTemplatesDir, c.templates); err != nil {
			logger.WithError(err).Error("failed to overload templates")
		}
	}

	return nil
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
