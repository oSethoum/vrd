package ent

import (
	"embed"
	"fmt"
	"os"
	"path"
	"vrd/config"
	"vrd/types"
	"vrd/utils"
)

//go:embed templates
var Templates embed.FS

type Engine struct {
	config *types.Config
	vuerd  *types.Vuerd
	files  []*types.File
	state  *State
}

// Create the Ent engine
func NewEngine() *Engine {
	c := config.Init()
	v := &types.Vuerd{}
	utils.ReadJson(c.Input, v)
	// show json
	// utils.WriteJson("vrd/output.vuerd.json", v)
	return &Engine{
		config: c,
		vuerd:  v,
		files:  []*types.File{},
	}
}

// Start the Ent engine
func (e *Engine) Start() {
	p := NewParser(*e.vuerd, *e.config)
	e.state = p.Start()

	// debug mode only
	if e.config.Debug {
		utils.WriteJson("vrd/output.json", e.state)
	}

	os.Exit(0)

	data := Data{
		Config: e.config,
		Nodes:  e.state.Nodes,
		Mixins: e.state.Mixins,
	}

	// nodes
	for _, node := range e.state.Nodes {
		data.Node = node
		e.files = append(e.files, &types.File{
			Path:   fmt.Sprintf("ent/schema/%s.go", node.Name),
			Buffer: e.parseTemplate("ent/schema/node.go.tmpl", data),
		})
	}

	// mixins
	for _, node := range e.state.Mixins {
		data.Mixin = node
		e.files = append(e.files, &types.File{
			Path:   fmt.Sprintf("ent/schema/%s.go", node.Name),
			Buffer: e.parseTemplate("ent/schema/mixin.go.tmpl", data),
		})
	}

	e.writeFiles()
}

func (e *Engine) writeFiles() {
	for _, file := range e.files {
		out := path.Join(e.config.Ent.Output, file.Path)
		utils.WriteFile(out, file.Buffer)
	}
}

func (e *Engine) parseTemplate(filename string, data interface{}) string {
	return utils.ParseTemplate(Templates, filename, data)
}
