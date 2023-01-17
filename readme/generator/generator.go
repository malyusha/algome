package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/malyusha/algome/logger"
	"github.com/malyusha/algome/util"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	readmeFilename = "README.md"
)

type Source struct {
	Name  string
	Stats Stats
}

type Stats struct {
	hideUnsolved bool
	problems     []Problem
	Solved       int64
}

func (s Stats) Total() int {
	return len(s.problems)
}

func (s *Stats) add(p Problem) {
	s.problems = append(s.problems, p)
	if p.IsSolved {
		s.Solved++
	}
}

func (s Stats) Problems() []Problem {
	if s.hideUnsolved {
		return util.Filter(s.problems, func(p Problem) bool { return p.IsSolved })
	}
	return s.problems
}

func newStats(hideUnsolved bool, problems []Problem) Stats {
	stats := Stats{
		hideUnsolved: hideUnsolved,
	}
	for _, p := range problems {
		stats.add(p)
	}
	return stats
}

func (s Stats) SolvedPercentString() string {
	if s.Solved == 0 {
		return "0"
	}

	res := float64(s.Solved*100) / float64(s.Total())

	if math.Mod(res, 1.0) == 0 {
		return strconv.Itoa(int(res))
	}

	return fmt.Sprintf("%.2f", res)
}

// ProblemsProvider - source of single source of algorithmic problems.
// E.g. - leetcode, hackerrank, etc.
type ProblemsProvider interface {
	// GetAllProblems returns list of all problems of algo Provider no matter solved they or not.
	GetAllProblems(ctx context.Context) ([]Problem, error)
}

type RenderFormatter interface {
	// FormatOutput formats list of problems and outputs bytes content to be written as w of
	// readme.
	FormatOutput(ctx context.Context, source Source) ([]byte, error)
}

type Solution struct {
	ProblemUID    string
	BasePath      string
	SolutionFiles []string
}

// StructureProvider is the source that knows directory & files structure of all solved problems.
type StructureProvider interface {
	GetSolvedProblems(ctx context.Context, sourceName string) ([]Solution, error)
	GetProblemUID(p Problem) string
}

// Generator is the main structure responsible for readme generation from different sources of problems.
type Generator struct {
	hideUnsolved      bool
	outputDir         string
	templates         *Templates
	structureProvider StructureProvider
	problemsProvider  []Provider
}

// NewGenerator returns new algo source readme generator.
func NewGenerator(
	outputDir string,
	templates *Templates,
	sp StructureProvider,
	pp []Provider,
	options ...Option,
) *Generator {
	gen := &Generator{
		outputDir:         outputDir,
		templates:         templates,
		structureProvider: sp,
		problemsProvider:  pp,
	}

	for _, o := range options {
		o(gen)
	}

	return gen
}

type GenerationResult struct {
	err error
}

// GenerateReadme ...
func (g *Generator) GenerateReadme(ctx context.Context) error {
	templateSources := make([]TemplateSource, 0, len(g.problemsProvider))
	for _, provider := range g.problemsProvider {
		if err := g.generateProviderReadme(ctx, provider); err != nil {
			return fmt.Errorf("error on source %s: %w", provider.name, err)
		}

		templateSources = append(templateSources, TemplateSource{
			Title: cases.Title(language.English, cases.Compact).String(provider.name),
			Path:  filepath.Join(g.outputDir, provider.name),
		})
	}

	indexFileWriter := newFileWriter(filepath.Join(g.outputDir, readmeFilename))
	if err := WriteReadme(indexFileWriter, g.templates.Index, templateSources); err != nil {
		return fmt.Errorf("failed to write index readme: %w", err)
	}

	return nil
}

func (g *Generator) generateProviderReadme(ctx context.Context, provider Provider) error {
	allProblems, err := provider.GetAllProblems(ctx)
	if err != nil {
		return fmt.Errorf("generateProviderReadme: %w", err)
	}

	solutions, err := g.structureProvider.GetSolvedProblems(ctx, provider.name)
	if err != nil {
		return fmt.Errorf("generateProviderReadme: %w", err)
	}

	// create map of unique ID of each problem and positions of corresponding problem inside
	// `allProblems` slice.
	problemsMap := make(map[string]int, len(allProblems))
	for ix := range allProblems {
		p := allProblems[ix]
		p.source = provider.name
		p.uid = g.structureProvider.GetProblemUID(p)
		problemsMap[p.uid] = ix
		allProblems[ix] = p
	}

	// filtering problems, keeping only solved
	for _, s := range solutions {
		ix := problemsMap[s.ProblemUID]
		problem := allProblems[ix]
		allProblems[ix] = problem.applySolution(s)
	}

	sort.Slice(allProblems, func(i, j int) bool { return allProblems[i].ID < allProblems[j].ID })
	stats := newStats(g.hideUnsolved, allProblems)

	stats.hideUnsolved = g.hideUnsolved
	source := Source{
		Name:  cases.Title(language.English, cases.Compact).String(provider.name),
		Stats: stats,
	}

	fw := newFileWriter(filepath.Join(g.outputDir, provider.name, readmeFilename))
	if err := WriteReadme(fw, g.templates.Source, source); err != nil {
		return fmt.Errorf("generateProviderReadme error: %w", err)
	}

	return nil
}

type Provider struct {
	name     string
	cacheDir string
	base     ProblemsProvider
}

func (p *Provider) GetAllProblems(ctx context.Context) ([]Problem, error) {
	problems, err := p.getProblemsFromCache()
	if err != nil || len(problems) == 0 {
		problems, err = p.base.GetAllProblems(ctx)
		if err != nil {
			return problems, err
		}
	}

	if err := p.writeProblemsToCache(problems); err != nil {
		logger.WithError(err).Error("failed to write problems into cache")
	}

	return problems, nil
}

func (p *Provider) writeProblemsToCache(problems []Problem) error {
	filename := fmt.Sprintf("_cache.%s.json", p.name)
	filePath := filepath.Join(p.cacheDir, filename)
	_, err := os.Stat(p.cacheDir)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(p.cacheDir, 0755); err != nil {
			return fmt.Errorf("failed to create cache dir: %w", err)
		}
	}

	cacheFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to marshal problems for cache: %w", err)
	}

	enc := json.NewEncoder(cacheFile)
	if err := enc.Encode(problems); err != nil {
		return fmt.Errorf("failed to write problems into cache file: %w", err)
	}

	return nil
}

func (p *Provider) getProblemsFromCache() ([]Problem, error) {
	filename := fmt.Sprintf("_cache.%s.json", p.name)
	filePath := filepath.Join(p.cacheDir, filename)

	contents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	dst := make([]Problem, 0)
	err = json.Unmarshal(contents, &dst)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal problems from cache: %w", err)
	}

	return dst, nil
}

func NewProvider(name string, cacheDir string, base ProblemsProvider) Provider {
	return Provider{
		name:     name,
		cacheDir: util.WithUserHomeDir(cacheDir),
		base:     base,
	}
}
