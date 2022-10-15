package ent

import "vrd/types"

type Data struct {
	Config *types.Config    `json:"config"`
	Nodes  []*Node          `json:"nodes"`
	Mixins map[string]*Node `json:"mixins"`
	Node   *Node            `json:"node"`
	Mixin  *Node            `json:"mixin"`
	types.Helper
}

type State struct {
	Nodes  []*Node          `json:"nodes"`
	Mixins map[string]*Node `json:"mixins"`
}

type Node struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Alias       string   `json:"alias"`
	TableName   string   `json:"tableName"`
	M2M         bool     `json:"m2m"`
	Comment     string   `json:"comment"`
	Fields      []Field  `json:"fields"`
	Edges       []Edge   `json:"edges"`
	Mixins      []string `json:"mixins"`
	Imports     []string `json:"imports"`
	Annotations []string `json:"annotations"`
}

type Edge struct {
	ID          string   `json:"id"`
	Node        string   `json:"node"`
	Name        string   `json:"name"`
	Reference   string   `json:"reference"`
	Options     []string `json:"options"`
	Direction   string   `json:"direction"` // To | From
	Annotations []string `json:"annotations"`
}

type Field struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	EnumValues  []string    `json:"enumValues"`
	Comment     string      `json:"comment"`
	Default     string      `json:"default"`
	Options     []string    `json:"options"`
	Skips       []string    `json:"skips"`
	Annotations []string    `json:"annotations"`
	Validation  *Validation `json:"validation"`
}

type Validation struct {
	Max         float32    `json:"max"`
	Min         float32    `json:"min"`
	MinLen      uint       `json:"minLen"`
	MaxLen      uint       `json:"maxLen"`
	Range       [2]float32 `json:"range"`
	Match       string     `json:"match"`
	Positive    bool       `json:"positive"`
	Negative    bool       `json:"negative"`
	NonNegative bool       `json:"nonNegative"`
	Immutable   bool       `json:"immutable"`
	Sensitive   bool       `json:"sensitive"`
	Optional    bool       `json:"optional"`
}

type EntSchema struct {
	Path        string `json:"path"`
	Schema      string `json:"schema"`
	Mixins      string `json:"mixins"`
	Fields      string `json:"fields"`
	Edges       string `json:"edges"`
	Annotations string `json:"annotations"`
	Policy      string `json:"policy"`
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
	Path          string `json:"path"`
	Head          string `json:"Head"`
	Queries       string `json:"queries"`
	Query         string `json:"query"`
	Create        string `json:"create"`
	Update        string `json:"update"`
	Delete        string `json:"delete"`
	Subscriptions string `json:"subscriptions"`
}
