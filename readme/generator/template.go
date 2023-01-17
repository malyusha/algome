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

var TplIndexMarkdown = `# All problems
{{- range . }}
[{{ .Title }}]({{ .Path }})
{{- end }}
`

var TplSourceMarkdown = `# {{ .Name }} problems
{{- with .Stats }}
{{- range $difficulty, $levelStats := groupByLevel . }}
{{- with $levelStats }}
## {{ title $difficulty }} - {{ .SolvedPercent }}% [{{ .Solved }} / {{ .Total }}]
{{- range .Problems }} 
	{{.ID}}. [{{.Title}}]({{.URL}}) 
{{- if .IsSolved }} 
{{- range .Solutions }}
		* [{{.Lang}}]({{.Filepath}})
{{- end }}
{{- end }}
{{- end }}
{{- end }}
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

func GroupByLevel(s Stats) map[string]Stats {
	groups := make(map[string]Stats)
	for _, p := range s.Problems {
		diff := p.Difficulty.String()
		stats := groups[diff]
		stats.add(p)
		groups[diff] = stats
	}

	return groups
}
