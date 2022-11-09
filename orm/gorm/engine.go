package gorm

import (
	"embed"
	"fmt"
	"path"
	"vrd/types"
	"vrd/utils"
)

//TODO: support DTO

//go:embed templates
var Templates embed.FS

type Engine struct {
	config  *types.Config
	vuerd   *types.Vuerd
	state   *State
	imports []string
	files   []*types.File
	h       *types.Helper
}

// Create the Ent engine
func NewEngine(c *types.Config) *Engine {
	v := &types.Vuerd{}
	utils.ReadJson(c.Input, v)
	return &Engine{
		config:  c,
		vuerd:   v,
		imports: []string{"gorm.io/gorm"},
		state:   &State{},
		files:   []*types.File{},
		h:       &types.Helper{},
	}
}

func (e *Engine) Start() {
	p := &Parser{e}
	p.Start()
	// debug mode only
	if e.config.Debug {
		utils.WriteJson("vrd/output.json", e.state)
	}

	data := Data{
		Package:   e.config.Gorm.Package,
		Database:  e.config.Database,
		Auth:      e.config.Gorm.Auth,
		Socket:    e.config.Gorm.Socket,
		Swagger:   e.config.Gorm.Swagger,
		Debug:     e.config.Debug,
		Imports:   e.imports,
		GormModel: e.config.Gorm.GormModel,
		Models:    e.state.Models,
	}

	files := []string{
		"config/config.go",
		"db/db.go",
		"models/models.go",
	}

	if e.config.Gorm.Fiber {
		files = append(files,
			"handlers/handlers.go",
			"handlers/types.go",
			"routes/routes.go",
			"main.go",
		)

		if data.Auth {
			files = append(files,
				"handlers/users.go",
				"handlers/login.go",
				"middlewares/auth.go",
			)
		}

		if data.Socket {
			files = append(files,
				"handlers/socket.go",
			)
		}
	}

	for _, file := range files {
		e.files = append(e.files,
			&types.File{
				Path:   file,
				Buffer: e.parseTemplate(file+".tmpl", data),
			},
		)
	}

	for _, m := range e.state.Models {
		data.Model = m
		if m.Name == "User" {
			continue
		}
		e.files = append(e.files,
			&types.File{
				Path:   fmt.Sprintf("handlers/%s.go", e.h.Snakes(m.Name)),
				Buffer: e.parseTemplate("handlers/model.go.tmpl", data),
			},
		)
	}

	if e.config.Gorm.Typescrpit {
		e.files = append(e.files, &types.File{
			Path:   "types/types.ts",
			Buffer: e.parseTemplate("ts.go.tmpl", data),
		})
	}

	e.files = append(e.files, &types.File{
		Path:   "vrd.sh",
		Buffer: e.parseTemplate("vrd.sh.go.tmpl", data),
	})

	e.writeFiles()
}

func (e *Engine) writeFiles() {
	for _, file := range e.files {
		out := path.Join(e.config.Gorm.Output, file.Path)
		utils.WriteFile(out, file.Buffer)
	}
}

func (e *Engine) parseTemplate(filename string, data interface{}) string {
	return utils.ParseTemplate(Templates, filename, data)
}
