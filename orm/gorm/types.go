package gorm

import "vrd/types"

type State struct {
	Models map[string]Model
}

type Data struct {
	Package   string
	Database  string
	Imports   []string
	Socket    bool
	Auth      bool
	Fiber     bool
	Swagger   bool
	Debug     bool
	GormModel bool
	Models    map[string]Model
	Model     Model
	types.Helper
}

type Model struct {
	Name       string
	Columns    map[string]Column
	GormModel  bool
	JsonFields []string
}

type Column struct {
	Name    string
	Type    string
	TsName  string
	TsType  string
	Options string
}
