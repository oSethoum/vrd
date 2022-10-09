package ent

import (
	"log"
	"os"
	"path"
	"strings"
	"vrd/config"
	"vrd/types"
	"vrd/utils"
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

func ParseTemplate(path string, data interface{}) string {
	return utils.ParseTemplate(Assets, path, data)
}
