package generator

type Option func(g *Generator)

// HideUnsolvedProblems returns option to be set for generator to hide unsolved problems.
func HideUnsolvedProblems() Option {
	return func(g *Generator) { g.hideUnsolved = true }
}
