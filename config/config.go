package config

import (
	"log"
	"os"
	"strings"
	"vrd/utils"
)

func Init() Config {
	var config Config
	_, err := os.Stat("vrd/vrd.config.json")

	if err != nil {
		config = Config{
			Input:  "vrd/db.vuerd.json",
			Output: "./",
			Ent: &Ent{
				Package: "app",
				Graphql: true,
				Echo:    true,
				Auth:    true,
				Privacy: true,
			},
		}
		utils.WriteJSON("vrd/db.vuerd.json", nil)
		utils.WriteJSON("vrd/vrd.config.json", config)
		println("vrd initialized successfully")
		os.Exit(0)
	} else {
		utils.ReadJSON("vrd/vrd.config.json", &config)
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
