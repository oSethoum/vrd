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

		println("Running: go mod init")
		err := exec.Command("go", "mod", "init", config.Ent.Package).Run()
		checkErr(err, false)
		if config.Ent.Privacy {
			config.Ent.PrivacyNode = false
		}

		println("Running: Ent Engine")
		ent.Engine(state, config)

		println("Running: go mod fmt")
		err = exec.Command("go", "fmt", "./...").Run()
		checkErr(err, false)

		println("Running: go mod tidy")
		err = exec.Command("go", "mod", "tidy").Run()
		checkErr(err, false)

		println("Running: go generate ./ent")
		err = exec.Command("go", "generate", "./ent").Run()
		checkErr(err, true)

		println("Running: go mod tidy")
		err = exec.Command("go", "mod", "tidy").Run()
		checkErr(err, false)

		if config.Ent.Privacy {
			config.Ent.PrivacyNode = true
			println("Running: Ent Engine With PrivacyNode")
			ent.Engine(state, config)
			println("Running: go generate ./ent")
			err = exec.Command("go", "generate", "./ent").Run()
			checkErr(err, false)
			println("Running: go mod tidy")
			err = exec.Command("go", "mod", "tidy").Run()
			checkErr(err, false)
		}
		println("Running: gqlgen")
		err = exec.Command("gqlgen").Run()
		checkErr(err, false)
	}
}

func checkErr(err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatal(err.Error())
		} else {
			log.Println("Warning", err.Error())
		}
	}
}
