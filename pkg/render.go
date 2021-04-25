package pkg

import (
	"embed"
	"html/template"
)

//go:embed *
var tmplFile embed.FS

type Result struct {
	output string
}

func (r *Result) Write(b []byte) (n int, err error) {
	r.output += string(b)
	return len(b), nil
}

func ProberRender(m interface{}) (result string, err error) {
	t, err := template.ParseFS(tmplFile, "prober.tmpl")
	if err != nil {
		return "", err
	}

	resultWriter := &Result{}
	if err := t.Execute(resultWriter, m); err != nil {
		return "", err
	}
	return resultWriter.output, nil
}
