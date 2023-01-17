package generator

import (
	"fmt"
	"text/template"
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
{{- with .Problems }}
	{{- range . }} 
		{{.ID}}. [{{.Title}}]({{.URL}}) 
		{{- if .IsSolved }} 
			{{- range .Solutions }}
			* [{{.Lang}}]({{.Filepath}}) 
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

	for name, text := range list {
		parsed, err := template.New(name).Parse(text)
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
