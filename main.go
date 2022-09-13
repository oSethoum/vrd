package main

import (
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
	}
}
