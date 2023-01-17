package generator

import (
	"fmt"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	TplSourceMarkdownName = "source_markdown"
	TplIndexMarkdownName  = "index"
)

var DefaultTemplates = map[string]string{
	TplIndexMarkdownName:  TplIndexMarkdown,
	TplSourceMarkdownName: TplSourceMarkdown,
}

var TplIndexMarkdown = `# Sources
{{- range . }}
## [{{ .Title }}]({{ .Path }})
{{- end }}

Readme generated using [algome](https://github.com/malyusha/algome).
`

var TplSourceMarkdown = `# {{ .Name }} problems
{{- with .Stats }}
{{- range groupByLevel . }}
<details>
	<summary>{{ title .Title }} - {{ .Stats.SolvedPercentString }}% [{{ .Stats.Solved }} / {{ .Stats.Total }}]</summary>
{{ range $i, $p := .Stats.Problems }} 
{{inc $i}}. [{{.ID}}. {{.Title}}]({{.URL}}) 
{{- if .IsSolved }} 
{{- range .Solutions }}
	* [{{.Lang}}]({{.Filepath}})
{{- end }}
{{- end -}}
{{- end }}
</details>
{{- end }}
{{- end }}
`

type Templates struct {
	Source, Index *template.Template
}

type TemplateSource struct {
	Title string
	Path  string
}

func LoadTemplates(out *Templates, list map[string]string) error {
	if out == nil {
		return fmt.Errorf("nil templates provided")
	}

	funcs := template.FuncMap{
		"groupByLevel": GroupByLevel,
		"title": func(s string) string {
			return cases.Title(language.English, cases.Compact).String(s)
		},
		"inc": func(i int) int {
			return i + 1
		},
	}
	for name, text := range list {
		parsed, err := template.New(name).Funcs(funcs).Parse(text)
		if err != nil {
			return fmt.Errorf("failed to parse template '%s': %w", name, err)
		}

		setTemplate := templateSetters[name]
		setTemplate(parsed, out)
	}

	return nil
}

type OverloadTemplate func(parsed *template.Template, t *Templates)

var templateSetters = map[string]OverloadTemplate{
	TplSourceMarkdownName: func(parsed *template.Template, t *Templates) { t.Source = parsed },
	TplIndexMarkdownName:  func(parsed *template.Template, t *Templates) { t.Index = parsed },
}

type ProblemsGroup struct {
	Title string
	Stats Stats
}

func GroupByLevel(s Stats) []ProblemsGroup {
	groupsMap := make(map[difficulty]Stats)
	for _, p := range s.problems {
		diff := p.Difficulty
		stats, ok := groupsMap[diff]
		if !ok {
			stats.hideUnsolved = s.hideUnsolved
		}
		stats.add(p)
		groupsMap[diff] = stats
	}

	order := []difficulty{DifficultyEasy, DifficultyMedium, DifficultyHard}
	out := make([]ProblemsGroup, 0, len(groupsMap))
	for _, diff := range order {
		if _, ok := groupsMap[diff]; !ok {
			continue
		}

		group := ProblemsGroup{
			Title: diff.String(),
			Stats: groupsMap[diff],
		}

		out = append(out, group)
	}
	return out
}
