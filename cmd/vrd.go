package main

import "vrd/orm/ent"

func main() {
	e := ent.NewEngine()
	e.Start()
}
