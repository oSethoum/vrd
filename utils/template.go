package utils

import (
	"bytes"
	"embed"
	"log"
	"strings"
	"text/template"
)

func ParseTemplate(assets embed.FS, fileName string, v interface{}) string {
	f, err := assets.ReadFile("templates/" + fileName)
	if err != nil {
		log.Fatalf("Engine: error reading file %s | %s", fileName, err.Error())
	}

	t, err := template.New(fileName).Parse(string(f))

	if err != nil {
		log.Fatalf("Engine: error parsing template %s | %s", fileName, err.Error())
	}

	out := bytes.Buffer{}

	err = t.Execute(&out, v)

	if err != nil {
		log.Fatalf("Engine: error executing template %s | %s", fileName, err.Error())
	}

	// correcting the characters
	str := strings.ReplaceAll(out.String(), "&#34;", "\"")
	str = strings.ReplaceAll(str, "&lt;", "<")
	return str
}
