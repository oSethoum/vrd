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
		ent.Engine(state, config)
		err := exec.Command("go", "fmt", "./...").Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}
