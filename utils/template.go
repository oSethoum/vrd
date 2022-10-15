package utils

import (
	"bytes"
	"embed"
	"strings"
	"text/template"
)

func ParseTemplate(assets embed.FS, fileName string, v interface{}) string {
	f, err := assets.ReadFile("templates/" + fileName)
	CatchError(Fatal, err)

	t, err := template.New(fileName).Parse(string(f))

	CatchError(Fatal, err)
	out := bytes.Buffer{}

	err = t.Execute(&out, v)
	CatchError(Fatal, err)

	// correcting the characters
	str := strings.ReplaceAll(out.String(), "&#34;", "\"")
	str = strings.ReplaceAll(str, "&lt;", "<")
	return str
}
