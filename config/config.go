package config

import (
	"embed"
	"log"
	"os"
	"strings"
	"vrd/utils"

	"gopkg.in/yaml.v3"
)

//go:embed templates
var Assets embed.FS

func Init() Config {
	var config Config

	_, err := os.Stat("vrd/db.vuerd.json")
	if err != nil {
		utils.WriteFile("vrd/db.vuerd.json", "")
	}
	_, err = os.Stat("vrd/vrd.config.yaml")

	if err != nil {
		println("vrd initialized successfully")
		utils.WriteFile("vrd/vrd.config.yaml", utils.ParseTemplate(Assets, "vrd.config.yaml.go.tmpl", nil))
		os.Exit(0)
	} else {
		b, _ := os.ReadFile("vrd/vrd.config.yaml")
		yaml.Unmarshal(b, &config)
		checkConfig(config)
	}

	return config
}

func checkConfig(config Config) {
	if len(config.Input) < 12 && !strings.HasSuffix(config.Input, "vuerd.json") {
		log.Fatalf("Config: error input doesn't follow pattern *.vuerd.json")
	}

	if config.Ent != nil {
		if config.Ent.Package == "" {
			log.Fatalf("Config: package must not be empty")
		}
	}
}
