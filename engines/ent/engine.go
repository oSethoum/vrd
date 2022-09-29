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
	st := Parse(state, config)
	utils.WriteJSON("vrd/output.json", st)

	files := []types.File{}

	entSchemas := []EntSchema{}
	for _, node := range st.Nodes {

		entSchema := EntSchema{
			Path:        fmt.Sprintf("ent/schema/%s.go", kace.Snake(node.Name)),
			Schema:      parseTemplate("schema.go.tmpl", node),
			Annotations: parseTemplate("annotations.go.tmpl", node),
		}

		if len(node.Fields) > 0 {
			entSchema.Fields = parseTemplate("fields.go.tmpl", node)
		}

		if len(node.Edges) > 0 {
			entSchema.Edges = parseTemplate("edges.go.tmpl", node)
		}

		if len(node.Mixins) > 0 {
			entSchema.Mixins = parseTemplate("mixins.go.tmpl", node)
		}

		entSchemas = append(entSchemas, entSchema)
	}

	entMixins := []EntMixin{}
	for _, mixin := range st.Mixins {
		entMixin := EntMixin{
			Path:   fmt.Sprintf("ent/schema/%s.go", kace.Snake(mixin.Name)),
			Schema: parseTemplate("mixin.schema.go.tmpl", mixin),
		}

		if len(mixin.Fields) > 0 {
			entMixin.Fields = parseTemplate("fields.go.tmpl", mixin)
		}

		if len(mixin.Edges) > 0 {
			entMixin.Edges = parseTemplate("edges.go.tmpl", mixin)
		}

		entMixins = append(entMixins, entMixin)
	}

	if config.Ent.Graphql {
		files = append(files,
			types.File{
				Path: "graph/resolvers/schema.resolvers.go",
				Buffer: parseTemplate("schema.resolvers.go.tmpl", SchemaData{
					Nodes:   st.Nodes,
					Package: config.Ent.Package,
				}),
			},
			types.File{
				Path:   "graph/resolvers/types.go",
				Buffer: parseTemplate("types.resolvers.go.tmpl", SchemaData{Nodes: st.Nodes, Package: config.Ent.Package}),
			},
			types.File{
				Path:   "graph/resolvers/notifiers.go",
				Buffer: parseTemplate("notifiers.resolvers.go.tmpl", SchemaData{Package: config.Ent.Package, Nodes: st.Nodes}),
			},
			types.File{
				Path:   "ent/generate.go",
				Buffer: parseTemplate("generate.go.tmpl", nil),
			},
			types.File{
				Path:   "graph/schemas/generated.graphqls",
				Buffer: parseTemplate("generated.go.tmpl", nil),
			},
			types.File{
				Path:   "ent/entc.go",
				Buffer: parseTemplate("entc.go.tmpl", nil),
			},
			types.File{
				Path:   "gqlgen.yml",
				Buffer: parseTemplate("gqlgen.go.tmpl", config.Ent.Package),
			},
			types.File{
				Path:   "graph/resolvers/resolver.go",
				Buffer: parseTemplate("resolver.go.tmpl", SchemaData{Package: config.Ent.Package, Nodes: st.Nodes}),
			},
			types.File{
				Path:   "handlers/handlers.go",
				Buffer: parseTemplate("handlers.go.tmpl", config.Ent.Package),
			},
		)

		gqlResolvers := []GQlResolver{}

		for _, node := range st.Nodes {
			gqlResolvers = append(gqlResolvers, GQlResolver{
				Path:          fmt.Sprintf("graph/resolvers/%s.resolvers.go", kace.Snake(node.Name)),
				Head:          parseTemplate("head.resolver.go.tmpl", config.Ent.Package),
				Query:         parseTemplate("query.resolvers.go.tmpl", node),
				Queries:       parseTemplate("queries.resolvers.go.tmpl", node),
				Create:        parseTemplate("create.resolvers.go.tmpl", node),
				Update:        parseTemplate("update.resolvers.go.tmpl", node),
				Delete:        parseTemplate("delete.resolvers.go.tmpl", node),
				Subscriptions: parseTemplate("subscriptions.resolvers.go.tmpl", node),
			})
			files = append(files, types.File{
				Path:   fmt.Sprintf("graph/schemas/%s.graphqls", kace.Snake(node.Name)),
				Buffer: parseTemplate("entity.graphqls.go.tmpl", node),
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
				types.File{
					Path:   "auth/types.go",
					Buffer: parseTemplate("auth.types.go.tmpl", nil),
				},
			)
		}

	}
	WriteSchemas(entSchemas, config)
	WriteMixins(entMixins, config)
	WriteFiles(files, config)
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

	str := strings.ReplaceAll(out.String(), "&#34;", "\"")
	str = strings.ReplaceAll(str, "&lt;", "<")
	return str
}
