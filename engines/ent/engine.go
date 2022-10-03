package ent

import (
	"embed"
	"fmt"

	"vrd/config"
	"vrd/types"
	"vrd/utils"
)

//go:embed templates
var Assets embed.FS

func Engine(state types.State, config config.Config) {
	st := Parse(state, config)
	st.Nodes = SortNodes(st.Nodes)
	utils.WriteJSON("vrd/output.json", st)

	files := []types.File{}
	data := TemplateData{
		Config: config,
		Nodes:  st.Nodes,
		Mixins: st.Mixins,
	}

	if config.Ent.Auth {
		files = append(files,
			types.File{
				Path:   "auth/login.go",
				Buffer: ParseTemplate("auth/login.go.tmpl", data),
			},
			types.File{
				Path:   "auth/middlewares.go",
				Buffer: ParseTemplate("auth/middlewares.go.tmpl", data),
			},
			types.File{
				Path:   "auth/types.go",
				Buffer: ParseTemplate("auth/types.go.tmpl", data),
			},
		)

		if config.Ent.Privacy {
			files = append(files, types.File{
				Path:   "auth/privacy.go",
				Buffer: ParseTemplate("auth/privacy.go.tmpl", data),
			})
		}
	}

	files = append(files,
		types.File{
			Path:   "db/db.go",
			Buffer: ParseTemplate("db/db.go.tmpl", data),
		},
		types.File{
			Path:   "server.go",
			Buffer: ParseTemplate("server.go.tmpl", data),
		},
		types.File{
			Path:   ".gitignore",
			Buffer: "vrd",
		},
	)

	if config.Ent.Graphql {
		files = append(files,
			types.File{
				Path:   "ent/entc.go",
				Buffer: ParseTemplate("ent/entc.go.tmpl", data),
			},
			types.File{
				Path:   "ent/generate.go",
				Buffer: ParseTemplate("ent/generate.go.tmpl", data),
			},
			types.File{
				Path:   "graph/resolvers/helpers.go",
				Buffer: ParseTemplate("graph/resolvers/helpers.go.tmpl", data),
			},
			types.File{
				Path:   "graph/resolvers/types.go",
				Buffer: ParseTemplate("graph/resolvers/types.go.tmpl", data),
			},
			types.File{
				Path:   "graph/resolvers/schema.resolvers.go",
				Buffer: ParseTemplate("graph/resolvers/schema.resolvers.go.tmpl", data),
			},
			types.File{
				Path:   "graph/resolvers/resolver.go",
				Buffer: ParseTemplate("graph/resolvers/resolver.go.tmpl", data),
			},
			types.File{
				Path:   "graph/schemas/types.graphqls",
				Buffer: ParseTemplate("graph/schemas/types.graphqls.go.tmpl", data),
			},
			types.File{
				Path:   "gqlgen.yml",
				Buffer: ParseTemplate("gqlgen.go.tmpl", data),
			},
			types.File{
				Path:   "routes/routes.go",
				Buffer: ParseTemplate("routes/routes.go.tmpl", data),
			},
			types.File{
				Path:   "handlers/handlers.go",
				Buffer: ParseTemplate("handlers/handlers.go.tmpl", data),
			},
			types.File{
				Path:   "graph/resolvers/utils.go",
				Buffer: ParseTemplate("graph/resolvers/utils.go.tmpl", data),
			},
		)

		if config.Ent.FileUpload {
			files = append(files,
				types.File{
					Path:   "graph/resolvers/upload.resolvers.go",
					Buffer: ParseTemplate("graph/resolvers/upload.resolvers.go.tmpl", nil),
				},
				types.File{
					Path:   "graph/schemas/upload.graphqls",
					Buffer: ParseTemplate("graph/schemas/upload.graphqls.go.tmpl", nil),
				},
			)
		}

		for index, node := range st.Nodes {
			data.Node = node
			data.Index = index
			files = append(files,
				types.File{
					Path:   fmt.Sprintf("graph/resolvers/%s.resolvers.go", node.Camel(node.Name)),
					Buffer: ParseTemplate("graph/resolvers/node.resolvers.go.tmpl", data),
				},
				types.File{
					Path:   fmt.Sprintf("graph/schemas/%s.graphqls", node.Camel(node.Name)),
					Buffer: ParseTemplate("graph/schemas/node.graphqls.go.tmpl", data),
				},
			)
		}
	}

	for _, node := range st.Nodes {
		data.Node = node
		files = append(files,
			types.File{
				Path:   fmt.Sprintf("ent/schema/%s.go", node.Camel(node.Name)),
				Buffer: ParseTemplate("ent/schema/node.go.tmpl", data),
			},
		)
	}

	for _, mixin := range st.Mixins {
		data.Mixin = mixin
		files = append(files,
			types.File{
				Path:   fmt.Sprintf("ent/schema/%s.go", mixin.Camel(mixin.Name)),
				Buffer: ParseTemplate("ent/schema/mixin.go.tmpl", data),
			},
		)
	}

	WriteFiles(files, config)
}
