package generator

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTemplates_Default(t *testing.T) {
	out := &Templates{}

	err := LoadTemplates(out, DefaultTemplates)
	assert.NoError(t, err)
}

func TestSourceTemplate(t *testing.T) {
	expectTemplateContents, err := ioutil.ReadFile("testdata/templates/source_markdown.md")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	problems := []Problem{
		{
			ID:         1,
			URL:        "https://problems/easy-one",
			Slug:       "easy-one",
			Title:      "Easy One",
			Difficulty: DifficultyEasy,
			source:     "test",
			uid:        "1",
			dir:        "test",
			Solutions: []SolutionFile{
				{
					Filename: "solution",
					Filepath: "test/solution.go",
					Lang:     LangGo,
				},
			},
			IsSolved: true,
		},
		{
			ID:         2,
			URL:        "https://problems/easy-two",
			Slug:       "easy-two",
			Title:      "Easy Two",
			Difficulty: DifficultyEasy,
			source:     "test",
			uid:        "2",
			dir:        "test",
			Solutions: []SolutionFile{
				{
					Filename: "solution",
					Filepath: "test/solution.go",
					Lang:     LangGo,
				},
			},
			IsSolved: true,
		},
		{
			ID:         3,
			URL:        "https://problems/hard-one",
			Slug:       "hard-one",
			Title:      "Hard One",
			Difficulty: DifficultyHard,
			source:     "test",
			uid:        "3",
			dir:        "test",
			IsSolved:   false,
		},
	}

	source := Source{
		Name:  "Test",
		Stats: newStats(problems),
	}

	templates := &Templates{}
	err = LoadTemplates(templates, DefaultTemplates)
	assert.NoError(t, err)
	buf := bytes.NewBuffer(nil)

	if !assert.NoError(t, templates.Source.Execute(buf, source)) {
		t.FailNow()
	}

	got := strings.Replace(buf.String(), "\t", "  ", -1)
	got = strings.Replace(got, " \n", "\n", -1)
	assert.Equal(t, string(expectTemplateContents), got)
}
