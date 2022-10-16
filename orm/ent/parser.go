package ent

import (
	"fmt"
	"strings"
	"vrd/types"
	"vrd/utils"
)

// Parser serve as a dependency injection
type Parser struct {
	config       *types.Config
	vuerd        *types.Vuerd
	state        *State
	currentNode  *Node
	currentField *Field
	currentEdge  *Edge
	h            types.Helper

	// use for comparison
	emptyGraphql types.Graphql
	emptyEnt     types.Ent
}

// Initialize new Parser with the needed dependencies
func NewParser(vuerd types.Vuerd, config types.Config) *Parser {
	// check for tables errors here
	return &Parser{
		state: &State{
			Nodes:  []*Node{},
			Mixins: make(map[string]*Node),
		},
		h:            types.Helper{},
		vuerd:        &vuerd,
		config:       &config,
		emptyGraphql: types.Graphql{},
		emptyEnt:     types.Ent{},
	}
}

func (p *Parser) Start() *State {

	if len(p.vuerd.TableState.Tables) == 0 {
		utils.CatchError(utils.Fatal, fmt.Errorf("`input.vuerd.json` has no tables"))
	}

	// Parsing mixins
	for i := 0; i < len(p.vuerd.TableState.Tables); i++ {
		if p.isMixin(p.vuerd.TableState.Tables[i].Name) {
			p.ParseTable(&p.vuerd.TableState.Tables[i])
		}
	}

	// Parsing tables
	for i := 0; i < len(p.vuerd.TableState.Tables); i++ {
		if !p.isMixin(p.vuerd.TableState.Tables[i].Name) {
			p.ParseTable(&p.vuerd.TableState.Tables[i])
		}
	}

	return p.state
}

// Parse Table
func (p *Parser) ParseTable(table *types.Table) {

	if p.h.Clean(table.Name) == "" {
		utils.CatchError(utils.Fatal, fmt.Errorf("table `%s` has no name", table.Id))
	}

	p.currentNode = &Node{
		ID:      table.Id,
		Name:    table.Name,
		Comment: table.Comment,
		Fields:  []Field{},
		Edges:   []Edge{},
		M2M:     strings.Contains(table.Comment, "-m2m"),
		Imports: []string{
			"\t\"entgo.io/ent\"",
		},
	}

	if p.currentNode.M2M {
		if len(table.Columns) < 3 {
			return
		} else {
			for _, column := range table.Columns {
				if column.Ui.Fk {
					p.parseColumn(&column)
				}
			}
		}
	}

	if p.config.Ent.Privacy && p.config.Ent.PrivacyNode {
		p.currentNode.Imports = append(p.currentNode.Imports, "\t\""+p.config.Ent.Package+"/ent/privacy\"")
		p.currentNode.Imports = append(p.currentNode.Imports, "\t\""+p.config.Ent.Package+"/auth\"")
	}

	if p.config.Ent.Graphql != p.emptyGraphql {
		p.currentNode.Imports = append(p.currentNode.Imports, "\t\"entgo.io/contrib/entgql\"")
	}

	for i := 0; i < len(table.Columns); i++ {
		if !table.Columns[i].Ui.Pk && !table.Columns[i].Ui.Fk && !table.Columns[i].Ui.Pfk {
			p.parseColumn(&table.Columns[i])
		}
	}

	for i := 0; i < len(p.vuerd.RelationshipState.Relationships); i++ {
		p.parseRelationship(&p.vuerd.RelationshipState.Relationships[i])
	}

	cms := p.h.Split(table.Comment, "|")

	if p.isMixin(table.Name) {
		p.currentNode.Imports = append(p.currentNode.Imports,
			"\t\"entgo.io/ent/schema/mixin\"",
		)

		for _, cm := range cms {
			opts := p.h.Split(cm, "=")

			if opts[0] == "nx" {
				p.currentNode.Alias = opts[1]
				p.state.Mixins[opts[1]] = p.currentNode
			} else {
				utils.CatchError(utils.Warninig, fmt.Errorf("table mixin `%s`, uknown option `%s`", p.currentNode.Name, opts[0]))
			}
		}

		if p.currentNode.Alias == "" {
			utils.CatchError(utils.Fatal, fmt.Errorf("table mixin `%s`, must have an alias eg: nx=time", p.currentNode.Name))
		}

		p.state.Mixins[p.currentNode.Alias] = p.currentNode

	} else {
		p.currentNode.Imports = append(p.currentNode.Imports,
			"\t\"entgo.io/ent/schema\"",
			"\t\"entgo.io/ent/dialect/entsql\"",
		)

		if p.config.Ent.Graphql != p.emptyGraphql {
			p.currentNode.Annotations = append(p.currentNode.Annotations,
				"entgql.RelayConnection()",
				"entgql.Mutation(entgql.MutationCreate(), entgql.MutationUpdate())",
			)
		}

		for _, cm := range cms {
			opts := p.h.Split(cm, "=")
			if len(opts) == 2 {
				if opts[0] == "mxs" {
					mxs := p.h.CleanSplit(opts[1], ",")
					for _, mx := range mxs {
						p.currentNode.Mixins = append(p.currentNode.Mixins, p.state.Mixins[mx].Name)
					}
				} else if opts[0] == "nm" {
					p.currentNode.TableName = opts[1]
				} else {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, uknwon option `%s`", opts[0], table.Name))
				}
			} else if len(opts) > 0 {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, option `%s` has no value", opts[0], table.Name))
			}
		}

		if p.currentNode.TableName == "" {
			p.currentNode.TableName = p.h.Snakes(p.currentNode.Name)
		}

		p.currentNode.Annotations = append(p.currentNode.Annotations,
			fmt.Sprintf("entsql.Annotation{Table: %s}", p.currentNode.TableName),
		)

		p.state.Nodes = append(p.state.Nodes, p.currentNode)
	}
}

func (p *Parser) parseColumn(column *types.Column) {

	if p.h.Clean(column.Name) == "" {
		utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, has no name", p.currentNode.Name, column.Id))
	}

	p.currentField = &Field{
		ID:          column.Id,
		Name:        column.Name,
		Comment:     column.Comment,
		Default:     column.Default,
		Options:     []string{},
		Annotations: []string{},
		Validation:  &Validation{},
	}

	p.parseColumnType(column)
	p.parseColumnOptions(column)
	p.currentNode.Fields = append(p.currentNode.Fields, *p.currentField)
}

func (p *Parser) parseRelationship(relationship *types.Relationship) {
	p.currentEdge = &Edge{
		ID: relationship.Id,
	}

	startTable := p.table(relationship.Start.TableId)
	endTable := p.table(relationship.End.TableId)
	endM2M := p.h.Contains(endTable.Comment, "-m2m")

	if endM2M && p.nonKeyColumnsCount(endTable) == 0 {
		return
	}

	endColumn := p.column(endTable, relationship.End.ColumnIds[0])
	if endM2M && startTable.Id == p.currentNode.ID && len(endTable.Columns) > 2 {
		p.parseM2MRelationship(startTable, endTable, endColumn)
	}

	cms := p.h.Split(endColumn.Comment, "|")
	names := []string{}
	name := ""

	for _, cm := range cms {
		if p.h.HasPreffix(cm, "nr=") {
			names = p.h.CleanSplit(cm, ",", "nr=(", ")")
			if len(names) != 2 {
				utils.CatchError(utils.Warninig, fmt.Errorf("table: `%s`, column: `%s`, nr=(start,end) must have two arguments", endTable.Name, endColumn.Name))
			}
		} else if p.h.HasPreffix(cm, "cr=") && p.currentNode.M2M {
			name = p.h.Clean("cr=")
		}
	}

	if startTable.Id == p.currentNode.ID {
		p.currentEdge.Direction = "To"
		p.currentEdge.Node = endTable.Name

		switch relationship.RelationshipType {
		case "ZeroN", "OneN":
			p.currentEdge.Node = endTable.Name

			if len(names) == 2 {
				p.currentEdge.Name = names[1]
			} else {
				p.currentEdge.Name = p.h.Snakes(endTable.Name)
			}

			if p.config.Ent.Graphql.RelayConnection {
				p.currentEdge.Annotations = append(p.currentEdge.Annotations, "entgql.RelayConnection()")
			}

		case "ZeroOne":
			if len(names) == 2 {
				p.currentEdge.Name = names[1]
			} else {
				p.currentEdge.Name = p.h.Snake(endTable.Name)
			}
			p.currentEdge.Options = append(p.currentEdge.Options, "Unique()")

		case "OneOnly":
			if len(names) == 2 {
				p.currentEdge.Name = names[1]
			} else {
				p.currentEdge.Name = p.h.Snake(endTable.Name)
			}
			p.currentEdge.Options = append(p.currentEdge.Options, "Unique()", "Required()")
		}

	} else if endTable.Id == p.currentNode.ID {
		if p.currentNode.M2M {
			p.currentEdge.Direction = "To"
			if name != "" {
				p.currentEdge.Name = name
			} else {
				p.currentEdge.Name = p.h.Snake(startTable.Name)
			}
			p.currentEdge.Options = append(p.currentEdge.Options,
				"Unique()",
				"Required()",
				fmt.Sprintf("Field(\"%s\")", endColumn.Name),
			)

		} else {
			p.currentEdge.Direction = "From"

			if len(names) == 2 {
				p.currentEdge.Name = names[0]
			} else {
				p.currentEdge.Name = p.h.Snake(startTable.Name)
			}

			switch relationship.RelationshipType {
			case "ZeroN", "ZeroOne":
				p.currentEdge.Options = append(p.currentEdge.Options, "Unique()")

			case "OneOnly", "OneN":
				p.currentEdge.Options = append(p.currentEdge.Options, "Unique()", "Required()")
			}
		}

	} else if startTable.Id == p.currentNode.ID && endTable.Id == p.currentNode.ID {
		p.currentEdge.Node = endTable.Name

		if len(names) == 2 {
			p.currentEdge.Name = names[0]
		} else {
			switch relationship.RelationshipType {
			case "ZeroN", "OneN":
				p.currentEdge.Name = p.h.Camels(endTable.Name)
			}
		}
	}
}

func (p *Parser) parseColumnType(column *types.Column) {
	column.DataType = p.h.Clean(column.DataType)

	if p.h.Contains(p.h.Pascal(column.DataType), "Enum") {
		p.currentField.Type = "Enum"
		p.currentField.EnumValues = p.h.CleanSplit(column.DataType, ",", "Enum", "(", ")")

		if p.config.Ent.Graphql != p.emptyGraphql {
			values := p.h.MultiplyArray(p.currentField.EnumValues, p.h.Pascal, p.h.UpperSnake)
			p.currentField.Options = append(p.currentField.Options, "NamedValues(\""+p.h.Join(values, "\", \"")+"\")")
		} else {
			values := p.h.MultiplyArray(p.currentField.EnumValues, p.h.Pascal)
			p.currentField.Options = append(p.currentField.Options, "Values(\""+p.h.Join(values, "\", \"")+"\")")
		}
	} else {
		t := types.VuerdTypes[p.h.Lower(column.DataType)]

		if t != "" {
			p.currentField.Type = EntTypesMap[t]
		} else {
			utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` supported", p.currentNode.Name, column.Name, column.DataType))
			p.currentField.Type = column.DataType
		}
	}
}

func (p *Parser) parseColumnDefault(value string, option string) {
	if p.currentField.Type == "Enum" {
		if !p.h.InArray(p.currentField.EnumValues, value) {
			utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, default value `%s` is not in enum", p.currentNode.Name, p.currentField.Name, value))
		} else if p.config.Ent.Graphql != p.emptyGraphql {
			p.currentField.Options = append(p.currentField.Options, option+"(\"", p.h.UpperSnake(value)+"\")")
		} else {
			p.currentField.Options = append(p.currentField.Options, option+"(\"", p.h.Pascal(value)+"\")")
		}
	} else if p.currentField.Type == "String" {
		p.currentField.Options = append(p.currentField.Options, option+`("`+value+`")`)
	} else {
		p.currentField.Options = append(p.currentField.Options, option+`(`+value+`)`)
	}

	if p.h.Contains(value, "time") && !p.h.InArray(p.currentNode.Imports, "\t\"time\"") {
		p.currentNode.Imports = append(p.currentNode.Imports, "\t\"time\"")
	}
}

func (p *Parser) parseColumnOptions(column *types.Column) {
	if p.config.Ent.Graphql != p.emptyGraphql {
		if p.h.InArray([]string{"created_at", "updated_at"}, column.Name) {
			p.currentField.Skips = append(p.currentField.Skips,
				"entgql.SkipMutationCreateInput",
				"entgql.SkipMutationUpdateInput",
			)
		}

		if p.h.InArray([]string{"password"}, column.Name) {
			p.currentField.Skips = append(p.currentField.Skips,
				"entgql.SkipWhereInput",
				"entgql.SkipType",
			)
		}
	}

	if column.Option.Unique {
		p.currentField.Options = append(p.currentField.Options, "Unique()")
	}

	if column.Option.AutoIncrement {
		p.currentField.Options = append(p.currentField.Options, "AutoIncrement()")
	}

	if !column.Option.NotNull {
		p.currentField.Options = append(p.currentField.Options, "Optional()", "Nillable()")
	}

	if len(column.Default) > 0 {
		p.parseColumnDefault(p.h.Clean(column.Default), "Default")
	}

	if p.config.Ent.Graphql != p.emptyGraphql && p.h.InArray(ComparableTypesMap, p.currentField.Type) {
		p.currentField.Options = append(p.currentField.Options, "entgql.OrderField(\""+p.h.UpperSnake(p.currentField.Name)+"\")")
	}

	cms := p.h.Split(column.Comment, "|")

	for _, cm := range cms {
		if p.h.HasPreffix(column.Comment, "upd=") {
			p.parseColumnDefault(p.h.Clean("upd="), "UpdateDefaut")
		} else if p.h.HasPreffix(cm, "minLen=") {

			if !p.h.InArray([]string{"String"}, p.currentField.Type) {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` doesn't support minLen", p.currentNode.Name, column.Name, p.currentField.Type))
				continue
			} else {
				value := p.h.Clean("minLen=")
				vv, ok := p.h.ValueOfType(value, "Uint")

				if !ok {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, value `%s` is not suitable", p.currentNode.Name, column.Name, value))
					continue
				}
				p.currentField.Validation.MinLen, _ = vv.(uint)
				p.currentField.Options = append(p.currentField.Options, fmt.Sprintf("MinLen(%s)", value))
			}
		} else if p.h.HasPreffix(cm, "maxLen=") {

			if !p.h.InArray([]string{"String"}, p.currentField.Type) {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` doesn't support minLen", p.currentNode.Name, column.Name, p.currentField.Type))
				continue
			} else {
				value := p.h.Clean("maxLen=")
				vv, ok := p.h.ValueOfType(value, p.currentField.Type)

				if !ok {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, value `%s` is not suitable form maxLen", p.currentNode.Name, p.currentField.Name, value))
					continue
				}
				p.currentField.Validation.MaxLen, _ = vv.(uint)
				p.currentField.Options = append(p.currentField.Options, fmt.Sprintf("MaxLen(%s)", value))
			}
		} else if p.h.HasPreffix(cm, "min=") {

			if !p.h.InArray([]string{"Int", "Uint", "Float"}, p.currentField.Type) {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` doesn't support min", p.currentNode.Name, p.currentField.Name, p.currentField.Type))
				continue
			} else {
				value := p.h.Clean("min=")
				v, ok := p.h.ValueOfType(value, p.currentField.Type)

				if !ok {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, value `%s` is not suitable", p.currentNode.Name, p.currentField.Name, value))
					continue
				}
				p.currentField.Validation.Min, _ = v.(float32)
				p.currentField.Options = append(p.currentField.Options, fmt.Sprintf("Min(%s)", value))
			}
		} else if p.h.HasPreffix(cm, "max=") {

			if !p.h.InArray([]string{"Int", "Uint", "Float"}, p.currentField.Type) {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` doesn't support max", p.currentNode.Name, p.currentField.Name, p.currentField.Type))
				continue
			} else {
				value := p.h.Clean("max=")
				vv, ok := p.h.ValueOfType(value, p.currentField.Type)

				if !ok {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, value `%s` is not suitable", p.currentNode.Name, p.currentField.Name, value))
					continue
				}
				p.currentField.Validation.Max, _ = vv.(float32)
				p.currentField.Options = append(p.currentField.Options, fmt.Sprintf("Max(%s)", value))
			}
		} else if p.h.HasPreffix(cm, "match=") {
			v := p.h.Clean("match=")

			if len(v) > 0 {
				p.currentField.Options = append(p.currentField.Options, "Match(\""+v+"\")")
			} else {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, option `%s` cannot be empty", p.currentNode.Name, p.currentField.Name, cm))
			}

		} else if p.h.HasPreffix(cm, "range=") {

			if !p.h.InArray([]string{"Float", "UIint", "Int"}, p.currentField.Type) {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` doesn't support range", p.currentNode.Name, p.currentField.Name, p.currentField.Type))
				continue
			}

			v := p.h.CleanSplit(cm, ",", "range=(", ")")
			if len(v) == 2 {
				i1, ok1 := p.h.ValueOfType(v[0], p.currentField.Type)
				i2, ok2 := p.h.ValueOfType(v[1], p.currentField.Type)

				v1, _ := i1.(float32)
				v2, _ := i2.(float32)

				p.currentField.Validation.Range = [2]float32{v1, v2}
				if !ok1 || !ok2 {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, range arguments must be of type `%s`", p.currentNode.Name, p.currentField.Name, p.currentField.Name))
					continue
				}

			} else {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, option `%s` must have 2 arguments", p.currentNode.Name, p.currentField.Name, cm))
			}
		} else if p.h.HasPreffix(cm, "skip=") {
			if p.config.Ent.Graphql != p.emptyGraphql {
				opts := p.h.CleanSplit(cm, ",", "skip=(", ")")
				for _, opt := range opts {
					if len(SkipMap[opt]) > 0 {
						if !p.h.InArray(p.currentField.Skips, SkipMap[opt]) {
							p.currentField.Skips = append(p.currentField.Skips, SkipMap[opt])
						}
					} else {
						utils.CatchError(utils.Warninig, fmt.Errorf("table `%s` column `%s`, option `%s` is not supported in skip", p.currentNode.Name, p.currentField.Name, opt))
					}
				}
				if len(p.currentField.Skips) > 0 {
					p.currentField.Annotations = append(p.currentField.Annotations, p.h.Join(p.currentField.Skips, "entgql.Skip(", ",", ")"))
				}
			} else {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s` column `%s`, skip only supported in graphql", p.currentNode.Name, p.currentField.Name))
			}
		} else if len(OptionsMap[cm]) > 0 {
			p.currentField.Options = append(p.currentField.Options, OptionsMap[cm])
		} else {
			utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, unknown option `%s`", p.currentNode.Name, p.currentField.Name, cm))
		}
	}
}

//debug: fmt.Printf("🐛 table `%s`, column `%s`, \n", p.currentNode.Name, p.currentField.Name)

// Check if table is a mixin
func (Parser) isMixin(name string) bool {
	return strings.Contains(strings.ToLower(name), "mixin")
}

// Find Table by ID
func (p *Parser) table(id string) *types.Table {
	for i := 0; i < len(p.vuerd.TableState.Tables); i++ {
		if p.vuerd.TableState.Tables[i].Id == id {
			return &p.vuerd.TableState.Tables[i]
		}
	}
	return nil
}

func (Parser) column(table *types.Table, id string) *types.Column {
	for j := 0; j < len(table.Columns); j++ {
		if table.Columns[j].Id == id {
			return &table.Columns[j]
		}
	}
	return nil
}

func (p *Parser) parseM2MRelationship(startTable, middleTable *types.Table, endColumn *types.Column) {
	nr := []string{}
	cms := p.h.Split(endColumn.Comment, "|")
	for _, cm := range cms {
		if p.h.Contains(cm, "nr=") {
			nr = p.h.CleanSplit(cm, ",", "nr=(", ")")
		}
	}

	if len(nr) > 0 {
		if nr[0] == "To" {
			p.currentEdge.Node = nr[1]
			p.currentEdge.Name = nr[2]
			p.currentEdge.Options = append(p.currentEdge.Options, fmt.Sprintf("Through(%s, %s.Type)", nr[3], middleTable.Name))
		} else if nr[0] == "From" {
			p.currentEdge.Node = nr[1]
			p.currentEdge.Name = nr[2]
			p.currentEdge.Options = append(p.currentEdge.Options, fmt.Sprintf("Ref(%s)", nr[3]))
			p.currentEdge.Options = append(p.currentEdge.Options, fmt.Sprintf("Through(%s, %s.Type)", nr[4], middleTable.Name))

		}
	} else {
		utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, must have argument nr=(To|From, Table, name, ref?, through)", middleTable.Name, endColumn.Name))
	}
}

func (Parser) nonKeyColumnsCount(table *types.Table) int {
	count := 0
	for _, column := range table.Columns {
		if !column.Ui.Fk || !column.Ui.Pfk || !column.Ui.Pk {
			count++
		}
	}
	return count
}
