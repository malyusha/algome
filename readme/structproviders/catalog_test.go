package structproviders

import (
	"context"
	"fmt"
	"testing"

	"github.com/malyusha/algome/readme/generator"
	"github.com/stretchr/testify/assert"
)

func TestCatalogStructProvider_GetSolvedProblems(t *testing.T) {
	cases := []struct {
		problemIdType ProblemIdType
		sourceName    string
		name          string
		expectOut     []generator.Solution
		expectError   error
	}{
		{
			name:          "no errors",
			sourceName:    "catalog",
			problemIdType: ProblemIdTypeId,
			expectOut: []generator.Solution{
				{
					ProblemUID:    "1",
					BasePath:      "testdata/catalog/1-100/1",
					SolutionFiles: []string{"main.test.go", "main.test.py"},
				},
				{
					ProblemUID:    "101",
					BasePath:      "testdata/catalog/100-200/101",
					SolutionFiles: []string{"main.test.go"},
				},
			},
		},
		{
			name:          "no errors for slugged if",
			problemIdType: ProblemIdTypeSlug,
			sourceName:    "catalog_slug_id",
			expectOut: []generator.Solution{
				{
					ProblemUID:    "test-slug",
					BasePath:      "testdata/catalog_slug_id/100-200/test-slug",
					SolutionFiles: []string{"main.test.go"},
				},
			},
		},
		{
			name:          "error invalid problem ID",
			problemIdType: ProblemIdTypeId,
			sourceName:    "catalog_slug_id",
			expectError:   fmt.Errorf("invalid solution UID: not convertable to int"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			is := assert.New(t)
			provider := NewCatalogStructProvider(c.problemIdType, "testdata")
			gotOut, err := provider.GetSolvedProblems(context.Background(), c.sourceName)
			if err != nil && c.expectError != nil {
				is.ErrorContains(err, c.expectError.Error())
			} else {
				if !is.NoError(err) {
					t.FailNow()
				}
				is.Equal(c.expectOut, gotOut)
			}
		})
	}
}
