package structproviders

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/malyusha/algome/readme/generator"
)

type ProblemIdType string

const (
	ProblemIdTypeId   ProblemIdType = "id"
	ProblemIdTypeSlug               = "slug"
)

// catalogStructProvider expects solved problems to be structured as catalogs.
// each catalog has numeric representation of range of problems.
// e.g.
//
// ┌ SOURCE_NAME <-- directory, containing single algo source solution files & sub directories.
// ├── 1-100
// │   ├── 1 <-- ID of the problem
// │   ├── 2
// │   ├── ...
// │   └── 100
// └── 100-200
//    ├── 100
//    ├── 101
//    ├── ...
//    └── 200
//
type catalogStructProvider struct {
	baseDir string
	idType  ProblemIdType
}

// NewCatalogStructProvider returns new catalog provider of solved problems.
func NewCatalogStructProvider(idType ProblemIdType, baseDir string) *catalogStructProvider {
	return &catalogStructProvider{
		baseDir: baseDir,
		idType:  idType,
	}
}

func (c *catalogStructProvider) GetSolvedProblems(ctx context.Context, sourceName string) ([]generator.Solution, error) {
	dir := filepath.Join(c.baseDir, sourceName)
	solutions, err := CatalogGetSolvedProblems(dir)
	if err != nil {
		return nil, err
	}

	for _, s := range solutions {
		if err := isValidCatalogProblemId(s.ProblemUID, c.idType); err != nil {
			return nil, fmt.Errorf("invalid solution UID: %w", err)
		}
	}

	return solutions, nil
}

// CatalogGetSolvedProblems returns solved problems list, ordered in ascending order
// by problem ID, from given path.
// It collects all files inside directory recursively and creates solution structs from them.
func CatalogGetSolvedProblems(path string) ([]generator.Solution, error) {
	// see expected structure of directories as described in comment of provider above.
	solvedProblemsDirs, err := collectFilesOfDirectory(path)
	if err != nil {
		return nil, fmt.Errorf("GetSolvedProblems error: %w", err)
	}

	solutions := make([]generator.Solution, 0, len(solvedProblemsDirs))
	for dir, files := range solvedProblemsDirs {
		solution, err := catalogExtractSolution(dir, files)
		if err != nil {
			return nil, fmt.Errorf("failed to create solution from directory: %w", err)
		}

		solutions = append(solutions, solution)
	}

	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].ProblemUID < solutions[j].ProblemUID
	})

	return solutions, nil
}

// catalogExtractSolution returns Solution struct from given base directory and its files.
// If the last part of directory is not convertable to int (thus, it's not and ID of a problem) then
// the error is returned.
func catalogExtractSolution(dir string, files []string) (out generator.Solution, err error) {
	pathParts := strings.Split(dir, string(os.PathSeparator))
	if len(pathParts) == 0 {
		return out, fmt.Errorf("not valid dir path given: '%s'", dir)
	}

	problemIdString := pathParts[len(pathParts)-1]

	out = generator.Solution{
		ProblemUID:    problemIdString,
		BasePath:      dir,
		SolutionFiles: files,
	}

	return
}

func isValidCatalogProblemId(id string, typ ProblemIdType) error {
	switch typ {
	default:
		return fmt.Errorf("unknown id type: %s", typ)
	case ProblemIdTypeId:
		// check if string is convertable to int, thus it's a real problem ID
		if _, err := strconv.Atoi(id); err != nil {
			return fmt.Errorf("not convertable to int: %w", err)
		}

	case ProblemIdTypeSlug:
		return nil
	}

	return nil
}

// GetProblemUID returns ID value of provided problem as string.
func (c *catalogStructProvider) GetProblemUID(p generator.Problem) string {
	switch c.idType {
	case ProblemIdTypeSlug:
		return p.Slug
	}

	return strconv.FormatInt(p.ID, 10)
}
