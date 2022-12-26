package generator

import (
	"path/filepath"
	"strings"
)

type difficulty int8

const (
	DifficultyEasy difficulty = iota
	DifficultyMedium
	DifficultyHard
)

type SolutionFile struct {
	Filename string
	Filepath string
	Lang     string
}

func (d difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "easy"
	case DifficultyMedium:
		return "medium"
	case DifficultyHard:
		return "hard"
	}

	return ""
}

type Problem struct {
	ID int64 `json:"id"`

	URL   string `json:"url"` // external URL of problem
	Slug  string `json:"slug"`
	Title string `json:"title"`

	Difficulty difficulty `json:"diff"`

	source string
	uid    string // unique problem identifier. any string by which we can determine exact problem
	dir    string

	Solutions []SolutionFile `json:"-"`
	IsSolved  bool           `json:"-"`
}

func (p *Problem) applySolution(s Solution) Problem {
	p.IsSolved = true
	s.BasePath = strings.TrimPrefix(s.BasePath, p.source+"/")

	p.dir = s.BasePath
	for _, fname := range s.SolutionFiles {
		p.addSolution(filepath.Join(s.BasePath, fname))
	}

	return *p
}

func (p *Problem) addSolution(path string) {
	supported, lang := langFromPath(path)
	if !supported {
		return
	}

	solution := SolutionFile{
		Filepath: path,
		Filename: filepath.Base(path),
		Lang:     lang,
	}

	p.Solutions = append(p.Solutions, solution)
}
