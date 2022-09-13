package ent

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"strings"

	"vrd/config"
	"vrd/types"
	"vrd/utils"

	"github.com/codemodus/kace"
)

//go:embed templates
var Assets embed.FS

func Engine(state types.State, config config.Config) {
	nodes := Parse(state, config)
	utils.WriteJSON("output.json", nodes)
	files := []types.File{}

	entSchemas := []EntSchema{}
	for _, node := range nodes {
		entSchemas = append(entSchemas, EntSchema{
			Path:        fmt.Sprintf("ent/schema/%s.go", kace.Snake(node.Name)),
			Schema:      parseTemplate("schema.go.tmpl", node),
			Fields:      parseTemplate("fields.go.tmpl", node),
			Edges:       parseTemplate("edges.go.tmpl", node),
			Annotations: parseTemplate("annotations.go.tmpl", node),
		})
	}

	// utils.WriteJSON("schemas.json", entSchemas)

	if config.Ent.Graphql {
		files = append(files,
			types.File{
				Path: "graph/resolvers/schema.resolvers.go",
				Buffer: parseTemplate("schema.resolvers.go.tmpl", SchemaData{
					Nodes:   nodes,
					Package: config.Ent.Package,
				}),
			},
			types.File{
				Path:   "ent/generate.go",
				Buffer: parseTemplate("generate.go.tmpl", nil),
			},
			types.File{
				Path:   "ent/entc.go",
				Buffer: parseTemplate("entc.go.tmpl", nil),
			},
			types.File{
				Path:   "gqlgen.yml",
				Buffer: parseTemplate("gqlgen.go.tmpl", config.Ent.Package),
			},
		)

		gqlResolvers := []GQlResolver{}

		for _, node := range nodes {
			gqlResolvers = append(gqlResolvers, GQlResolver{
				Path:    fmt.Sprintf("graph/resolvers/%s.resolver.go", kace.Snake(node.Name)),
				Query:   parseTemplate("query.resolvers.go.tmpl", node.Name),
				Queries: parseTemplate("queries.resolvers.go.tmpl", QueriesData{Name: node.Name}),
				Create:  parseTemplate("create.resolvers.go.tmpl", node.Name),
				Update:  parseTemplate("update.resolvers.go.tmpl", node.Name),
				Delete:  parseTemplate("delete.resolvers.go.tmpl", node.Name),
			})
			files = append(files, types.File{
				Path:   fmt.Sprintf("graph/schemas/%s.graphqls", kace.Snake(node.Name)),
				Buffer: parseTemplate("entity.graphqls.go.tmpl", node.Name),
			})
		}
		files = append(files,
			types.File{
				Path:   "db/db.go",
				Buffer: parseTemplate("db.go.tmpl", config.Ent.Package),
			},
			types.File{
				Path:   "server.go",
				Buffer: parseTemplate("server.go.tmpl", config.Ent.Package),
			},
		)

		// utils.WriteJSON("resolvers.json", gqlResolvers)
		WriteResolvers(gqlResolvers, config)

		if config.Ent.Echo {
			files = append(files,
				types.File{
					Path:   "routes/routes.go",
					Buffer: parseTemplate("routes.go.tmpl", config.Ent.Package),
				},
				types.File{
					Path:   "auth/auth.go",
					Buffer: parseTemplate("auth.go.tmpl", config.Ent.Package),
				},
			)
		}

	}
	WriteSchemas(entSchemas, config)
	WriteFiles(files, config)
}

// Types
type SchemaData struct {
	types.Helper
	Nodes   []Node
	Package string
}
type QueriesData struct {
	types.Helper
	Name string
}

// Helpers
func parseTemplate(fileName string, v interface{}) string {
	f, err := Assets.ReadFile("templates/" + fileName)
	if err != nil {
		log.Fatalf("Engine: error reading file %s", fileName)
	}

	t, err := template.New(fileName).Parse(string(f))

	if err != nil {
		log.Fatalf("Engine: error parsing template %s", fileName)
	}

	out := bytes.Buffer{}

	err = t.Execute(&out, v)

	if err != nil {
		log.Fatalf("Engine: error executing template %s", fileName)
	}

	return strings.ReplaceAll(out.String(), "&#34;", "\"")
}
