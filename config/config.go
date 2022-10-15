package config

import (
	"embed"
	"os"
	"path"
	"vrd/types"
	"vrd/utils"

	"gopkg.in/yaml.v3"
)

//go:embed templates
var Templates embed.FS

func Init() *types.Config {
	var config *types.Config
	cwd, _ := os.Getwd()
	_, err := os.Stat(path.Join(cwd, "vrd/vrd.config.yaml"))
	if err != nil {
		// initialize config
		buffer := utils.ParseTemplate(Templates, "vrd.config.yaml.go.tmpl", nil)
		utils.WriteFile("vrd/vrd.config.yaml", buffer)
		_, err := os.Stat(path.Join(cwd, "vrd/input.vuerd.json"))
		if err != nil {
			utils.WriteFile("vrd/input.vuerd.json", "")
		}
		println("🟢 vrd init")
		os.Exit(0)
	} else {
		// read the config
		config = &types.Config{}
		buffer := utils.ReadFile("vrd/vrd.config.yaml")
		err := yaml.Unmarshal([]byte(buffer), config)
		utils.CatchError(utils.Fatal, err)
	}
	return config
}
