package ent

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path"
	"strings"
	"vrd/config"
	"vrd/types"
)

func WriteFile(file types.File, c config.Config) {
	cwd, _ := os.Getwd()
	outPath := path.Join(cwd, c.Output, file.Path)
	os.MkdirAll(path.Dir(outPath), 0666)
	err := os.WriteFile(outPath, []byte(file.Buffer), 0666)
	if err != nil {
		log.Fatalf("Writing file %s: %v", file.Path, err)
	}
}
func WriteFiles(files []types.File, c config.Config) {
	for _, file := range files {
		WriteFile(file, c)
	}
}

func ParseTemplate(fileName string, v interface{}) string {
	f, err := Assets.ReadFile("templates/" + fileName)
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

func SortNodes(nodes []Node) []Node {
	for i := 0; i < len(nodes)-1; i++ {
		for j := i; j < len(nodes); j++ {
			if strings.Compare(nodes[i].Name, nodes[j].Name) == 1 {
				node := nodes[i]
				nodes[i] = nodes[j]
				nodes[j] = node
			}
		}
	}

	return nodes
}
