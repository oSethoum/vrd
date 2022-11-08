package main

import (
	"vrd/config"
	"vrd/orm/ent"
	"vrd/orm/gorm"
)

func main() {
	c := config.Init()

	if c.Ent != nil {
		e := ent.NewEngine(c)
		e.Start()
	}

	if c.Gorm != nil {
		e := gorm.NewEngine(c)
		e.Start()
	}

}
