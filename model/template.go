package model

import (
	"bytes"
	"os"
	"text/template"
)

func ApplyTemplate(s string) (string, error) {
	t, err := template.New("config").Funcs(template.FuncMap{
		"env": func(name string) string {
			return os.Getenv(name)
		},
		"foo": func() string {
			return "aaa"
		},
	}).Parse(s)
	if err != nil {
		return "", err
	}
	data := map[string]interface{}{}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
