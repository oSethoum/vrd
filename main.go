package main

import (
	"log"
	"os/exec"
	"vrd/config"
	"vrd/engines/ent"
	"vrd/types"
	"vrd/utils"
)

func main() {
	config := config.Init()
	var state types.State
	utils.ReadJSON(config.Input, &state)
	if config.Ent != nil {
		err := exec.Command("go", "mod", "init").Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
		ent.Engine(state, config)
		err = exec.Command("go", "fmt", "./...").Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = exec.Command("go", "mod", "tidy").Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = exec.Command("go", "generate", "./...").Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = exec.Command("gqlgen").Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
