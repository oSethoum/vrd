package ent

import (
	"vrd/types"
)

type State struct {
	Nodes  []Node           `json:"nodes"`
	Mixins map[string]Mixin `json:"mixins"`
}

type Mixin struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TableName   string   `json:"tableName"`
	Alias       string   `json:"alias"`
	Comment     string   `json:"comment"`
	Fields      []Field  `json:"fields"`
	Edges       []Edge   `json:"edges"`
	Imports     []string `json:"imports"`
	Annotations []string `json:"annotations"`
	types.Helper
}

type Node struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	TableName   string   `json:"tableName"`
	Comment     string   `json:"comment"`
	Fields      []Field  `json:"fields"`
	Edges       []Edge   `json:"edges"`
	Mixins      []string `json:"mixins"`
	Imports     []string `json:"imports"`
	Annotations []string `json:"annotations"`
	types.Helper
}

type Edge struct {
	ID          string   `json:"id"`
	Node        string   `json:"node"`
	Name        string   `json:"name"`
	Reference   string   `json:"reference"`
	Options     []string `json:"options"`
	Direction   string   `json:"direction"` // To | From
	Annotations []string `json:"annotations"`
	types.Helper
}

type Field struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	EnumValues  []string `json:"enumValues"`
	Comment     string   `json:"comment"`
	Default     string   `json:"default"`
	Options     []string `json:"options"`
	Annotations []string `json:"annotations"`
	types.Helper
}

type EntSchema struct {
	Path        string `json:"path"`
	Schema      string `json:"schema"`
	Mixins      string `json:"mixins"`
	Fields      string `json:"fields"`
	Edges       string `json:"edges"`
	Annotations string `json:"annotations"`
}

type EntMixin struct {
	Path        string `json:"path"`
	Schema      string `json:"schema"`
	Mixins      string `json:"mixins"`
	Fields      string `json:"fields"`
	Edges       string `json:"edges"`
	Annotations string `json:"annotations"`
}

type GQlResolver struct {
	Path    string `json:"path"`
	Head    string `json:"Head"`
	Queries string `json:"queries"`
	Query   string `json:"query"`
	Create  string `json:"create"`
	Update  string `json:"update"`
	Delete  string `json:"delete"`
}

var EntTypes = map[string]string{
	"int":      "Int",
	"long":     "Int64",
	"float":    "Float",
	"uuid":     "UUID",
	"double":   "Float64",
	"decimal":  "Int",
	"boolean":  "Bool",
	"string":   "String",
	"lob":      "String",
	"date":     "Time",
	"json":     "Json",
	"datetime": "Time",
	"time":     "Time",
}

var ComparableTypes = []string{
	"Int",
	"Int64",
	"Float",
	"UUID",
	"Float64",
	"Int",
	"Bool",
	"String",
	"String",
	"Time",
	"UUID",
}
