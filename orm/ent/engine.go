package ent

import (
	"embed"
	"fmt"
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
	h      *types.Helper
}

// Create the Ent engine
func NewEngine() *Engine {
	c := config.Init()
	v := &types.Vuerd{}
	utils.ReadJson(c.Input, v)
	return &Engine{
		config: c,
		vuerd:  v,
		files:  []*types.File{},
		h:      &types.Helper{},
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

	// os.Exit(0)

	data := Data{
		Config: e.config,
		Nodes:  e.state.Nodes,
		Mixins: e.state.Mixins,
	}

	// nodes
	for _, node := range e.state.Nodes {
		data.Node = node
		e.files = append(e.files, &types.File{
			Path:   fmt.Sprintf("ent/schema/%s.go", e.h.Snake(node.Name)),
			Buffer: e.parseTemplate("ent/schema/node.go.tmpl", data),
		})

		if e.config.Ent.Graphql != nil {
			e.files = append(e.files, &types.File{
				Path:   fmt.Sprintf("graph/schemas/%s.graphqls", e.h.Snake(node.Name)),
				Buffer: e.parseTemplate("graph/schemas/node.graphqls.go.tmpl", data),
			})
			e.files = append(e.files, &types.File{
				Path:   fmt.Sprintf("graph/resolvers/%s.resolvers.go", e.h.Snake(node.Name)),
				Buffer: e.parseTemplate("graph/resolvers/node.resolvers.go.tmpl", data),
			})
		}
	}

	// mixins
	for _, node := range e.state.Mixins {
		data.Mixin = node
		e.files = append(e.files, &types.File{
			Path:   fmt.Sprintf("ent/schema/%s.go", e.h.Snake(node.Name)),
			Buffer: e.parseTemplate("ent/schema/mixin.go.tmpl", data),
		})
	}

	files := []string{"db/db.go", "ent/generate.go"}

	if e.config.Ent.Echo {
		files = append(files, "server.go", "routes/routes.go", "handlers/handlers.go")
	}

	if e.config.Ent.Graphql != nil {
		files = append(files,
			"graph/schemas/types.graphqls", "gqlgen.yml", "graph/resolvers/resolver.go",
			"graph/resolvers/schema.resolvers.go", "ent/entc.go",
		)

		if e.config.Ent.Graphql.FileUpload {
			files = append(files, "graph/schemas/upload.graphqls", "graph/resolvers/upload.resolvers.go")
		}

		if e.config.Ent.Graphql.Subscription {
			files = append(files,
				"graph/resolvers/types.go", "graph/resolvers/utils.go", "graph/resolvers/notifiers.go",
			)
		}

	}

	if e.config.Ent.Auth {
		files = append(files,
			"auth/middlewares.go",
			"auth/types.go",
			"auth/login.go",
		)

		if e.config.Ent.Privacy {
			files = append(files, "auth/privacy.go")
		}
	}

	for _, file := range files {
		filename := file
		if p.h.HasSuffix(filename, ".go") {
			filename += ".tmpl"
		} else {
			filename += ".go.tmpl"
		}
		e.files = append(e.files,
			&types.File{
				Path:   file,
				Buffer: e.parseTemplate(filename, data),
			},
		)
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
