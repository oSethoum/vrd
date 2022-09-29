package ent

import (
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

func WriteResolvers(resolvers []GQlResolver, c config.Config) {
	cwd, _ := os.Getwd()

	for _, r := range resolvers {
		_, _ = os.Stat(path.Join(cwd, r.Path))
		//if os.IsNotExist(err) {
		WriteFile(types.File{
			Path: r.Path,
			Buffer: strings.Join([]string{
				r.Head,
				r.Query,
				r.Create,
				r.Update,
				r.Delete,
				r.Subscriptions,
			}, "\n\n"),
		}, c)
		//}
	}
}

func WriteSchemas(schemas []EntSchema, c config.Config) {
	cwd, _ := os.Getwd()

	for _, s := range schemas {
		_, _ = os.Stat(path.Join(cwd, s.Path))
		//if os.IsNotExist(err) {
		WriteFile(types.File{
			Path: s.Path,
			Buffer: strings.Join([]string{
				s.Schema,
				s.Mixins,
				s.Fields,
				s.Edges,
				s.Annotations,
				s.Policy,
			}, "\n\n"),
		}, c)
		//}
	}
}

func WriteMixins(mixins []EntMixin, c config.Config) {
	cwd, _ := os.Getwd()

	for _, s := range mixins {
		_, _ = os.Stat(path.Join(cwd, s.Path))
		//if os.IsNotExist(err) {
		WriteFile(types.File{
			Path: s.Path,
			Buffer: strings.Join([]string{
				s.Schema,
				s.Fields,
				s.Edges,
				s.Annotations,
			}, "\n\n"),
		}, c)
		//}
	}
}
