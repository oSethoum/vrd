package ent

import (
	"fmt"
	"regexp"
	"strings"
	"vrd/types"
	"vrd/utils"
)

type Parser struct {
	config *types.Config
	vuerd  *types.Vuerd
	state  *State
	cNode  *Node
	cField *Field
	h      types.Helper
}

func NewParser(vuerd types.Vuerd, config types.Config) *Parser {
	return &Parser{
		state: &State{
			Nodes:  []*Node{},
			Mixins: make(map[string]*Node),
		},
		h:      types.Helper{},
		vuerd:  &vuerd,
		config: &config,
	}
}

func (p *Parser) Start() *State {

	if len(p.vuerd.TableState.Tables) == 0 {
		utils.CatchError(utils.Fatal, fmt.Errorf("`input.vuerd.json` has no tables"))
	}

	for i := 0; i < len(p.vuerd.TableState.Tables); i++ {
		if p.isMixin(p.vuerd.TableState.Tables[i].Name) {
			p.ParseTable(&p.vuerd.TableState.Tables[i])
		}
	}

	for i := 0; i < len(p.vuerd.TableState.Tables); i++ {
		if !p.isMixin(p.vuerd.TableState.Tables[i].Name) {
			p.ParseTable(&p.vuerd.TableState.Tables[i])
		}
	}

	for i := 0; i < len(p.vuerd.RelationshipState.Relationships); i++ {
		p.parseRelationship(&p.vuerd.RelationshipState.Relationships[i])
	}

	return p.state
}

func (p *Parser) ParseTable(table *types.Table) {
	if p.h.Clean(table.Name) == "" {
		utils.CatchError(utils.Fatal, fmt.Errorf("table `%s` has no name", table.Id))
	}

	p.cNode = &Node{
		ID:      table.Id,
		Name:    p.h.Clean(table.Name),
		Comment: p.h.Clean(table.Comment),
		Fields:  []Field{},
		Edges:   []Edge{},
		M2M:     p.h.Contains(table.Comment, "-m2m"),
		Imports: []string{
			"\t\"entgo.io/ent\"",
		},
	}

	if p.cNode.M2M {
		p.validateM2MTable(table)
		if p.nonKeyColumnsCount(table) == 0 {
			return
		} else {
			for _, column := range table.Columns {
				if column.Ui.Fk && len(p.commentOptionValues(column.Comment, "dr=")) > 0 {
					p.parseColumn(&column)
				}
			}
		}
	}

	if p.config.Ent.Privacy && p.config.Ent.PrivacyNode {
		p.cNode.Imports = append(p.cNode.Imports, "\t\""+p.config.Ent.Package+"/ent/privacy\"")
		p.cNode.Imports = append(p.cNode.Imports, "\t\""+p.config.Ent.Package+"/auth\"")
	}

	if p.config.Ent.Graphql != nil {
		p.cNode.Imports = append(p.cNode.Imports, "\t\"entgo.io/contrib/entgql\"")
	}

	for i := 0; i < len(table.Columns); i++ {
		if !table.Columns[i].Ui.Pk && !table.Columns[i].Ui.Fk && !table.Columns[i].Ui.Pfk {
			p.parseColumn(&table.Columns[i])
		}
	}

	cms := p.h.Split(table.Comment, "|")

	if p.isMixin(table.Name) {
		p.cNode.Imports = append(p.cNode.Imports,
			"\t\"entgo.io/ent/schema/mixin\"",
		)

		for _, cm := range cms {
			opts := p.h.Split(cm, "=")

			if opts[0] == "nx" {
				p.cNode.Alias = opts[1]
				p.state.Mixins[opts[1]] = p.cNode
			} else {
				utils.CatchError(utils.Warninig, fmt.Errorf("table mixin `%s`, uknown option `%s`", p.cNode.Name, opts[0]))
			}
		}

		if p.cNode.Alias == "" {
			utils.CatchError(utils.Fatal, fmt.Errorf("table mixin `%s`, must have an alias eg: nx=time", p.cNode.Name))
		}

		p.state.Mixins[p.cNode.Alias] = p.cNode

	} else {
		p.cNode.Imports = append(p.cNode.Imports,
			"\t\"entgo.io/ent/schema\"",
			"\t\"entgo.io/ent/dialect/entsql\"",
		)

		if p.config.Ent.Graphql != nil {
			p.cNode.Annotations = append(p.cNode.Annotations,
				fmt.Sprintf("entgql.QueryField(\"%s\")", p.h.MCamels(p.cNode.Name)),
				"entgql.RelayConnection()",
				"entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate())",
			)
		}

		for _, cm := range cms {
			opts := p.h.Split(cm, "=")
			if len(opts) == 2 {
				if opts[0] == "mxs" {
					mxs := p.h.CleanSplit(opts[1], ",")
					for _, mx := range mxs {
						p.cNode.Mixins = append(p.cNode.Mixins, p.state.Mixins[mx].Name)
					}
				} else if opts[0] == "nm" {
					p.cNode.TableName = opts[1]
				} else {
					utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, uknwon option `%s`", opts[0], table.Name))
				}
			} else if len(opts) > 0 && !p.h.InArray(opts, "-m2m") {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, option `%s` has no value", table.Name, opts[0]))
			}
		}

		if p.cNode.TableName == "" {
			p.cNode.TableName = p.h.Snakes(p.cNode.Name)
		}

		p.cNode.Annotations = append(p.cNode.Annotations,
			fmt.Sprintf("entsql.Annotation{Table: \"%s\"}", p.cNode.TableName),
		)

		p.state.Nodes = append(p.state.Nodes, p.cNode)
	}
}

func (p *Parser) parseColumn(column *types.Column) {
	if p.h.Clean(column.Name) == "" {
		utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, has no name", p.cNode.Name, column.Id))
	}

	p.addImports("\t\"entgo.io/ent/schema/field\"")

	p.cField = &Field{
		ID:          column.Id,
		Name:        p.h.Clean(column.Name),
		Comment:     p.h.Clean(column.Comment),
		Default:     p.h.Clean(column.Default),
		Options:     []string{},
		Annotations: []string{},
	}

	p.parseColumnType(column)

	if !column.Ui.Fk {
		p.parseColumnOptions(column)

		if len(p.cField.Skips) > 0 {
			p.cField.Annotations = append(p.cField.Annotations, p.h.Join(p.cField.Skips, "entgql.Skip(", ",", ")"))
		}

		if len(p.cField.Annotations) > 0 {
			p.cField.Options = append(p.cField.Options, p.h.Join(p.cField.Annotations, "Annotations(", ",", ")"))
		}
	}

	p.cNode.Fields = append(p.cNode.Fields, *p.cField)
}

func (p *Parser) parseRelationship(relationship *types.Relationship) {
	startNode := p.node(relationship.Start.TableId)
	endNode := p.node(relationship.End.TableId)
	endColumn := p.column(p.table(relationship.End.TableId), relationship.End.ColumnIds[0])

	p.addImports("\t\"entgo.io/ent/schema/edge\"", startNode, endNode)
	// investigate this
	if p.h.Contains(endNode.Comment, "-m2m") {
		if len(endNode.Fields) > 0 && len(p.commentOptionValues(endColumn.Comment, "dr=")) > 0 {
			relationship.RelationshipType = "NN"
		} else {
			relationship.RelationshipType = "M2M"
		}
	}

	startEdge := &Edge{
		ID:        relationship.Id,
		Node:      endNode.Name,
		Name:      p.h.Snake(endNode.Name),
		Direction: "To",
	}

	endEdge := &Edge{
		ID:        relationship.Id,
		Node:      startNode.Name,
		Name:      p.h.Snake(startNode.Name),
		Direction: "From",
	}

	names := p.commentOptionValues(endColumn.Comment, "nr=")

	if len(names) == 2 {
		startEdge.Name = names[0]
		endEdge.Name = names[1]
	} else if len(names) == 1 {
		startEdge.Name = names[0]
	}

	if p.h.Contains(endColumn.Comment, "-m2m") {
		if p.h.Contains(endColumn.Comment, "-bi") {
			if startNode.ID == endNode.ID {
				relationship.RelationshipType = "M2MBI"
			} else {
				utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, -bi only allowed in same type edge", endNode.Name, endColumn.Name))
			}
		} else {
			relationship.RelationshipType = "M2M"
		}
	} else {
		if p.h.Contains(endColumn.Comment, "-bi") {
			if startNode.ID == endNode.ID {
				relationship.RelationshipType = "O2OBI"
			} else {
				utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, -bi only allowed in same type edge", endNode.Name, endColumn.Name))
			}
		}
	}

	switch relationship.RelationshipType {
	case "ZeroOne":
		startEdge.Options = []string{"Unique()"}
		endEdge.Options = []string{"Unique()"}
	case "ZeroN":
		endEdge.Options = []string{"Unique()"}
	case "OneOnly":
		startEdge.Options = []string{"Unique()"}
		endEdge.Options = []string{"Unique()", "Required()"}
	case "OneN":
		endEdge.Options = []string{"Unique()"}
	case "NN":
		direction := p.commentOptionValues(endColumn.Comment, "dr=")
		through := p.commentOptionValues(endColumn.Comment, "th=")

		startEdge.Direction = direction[0]
		startEdge.Node = direction[1]

		endEdge.Direction = "To"
		endEdge.Options = append(endEdge.Options, fmt.Sprintf("Field(\"%s\")", endColumn.Name))

		if len(through) == 0 {
			through = append(through, p.h.Snakes(endNode.Name))
		}

		startEdge.Options = append(startEdge.Options, fmt.Sprintf("Through(\"%s\", %s.Type)", through[0], endNode.Name))

	case "M2MBI":
		startNode.Edges = append(startNode.Edges, *startEdge)
	case "O2OBI":
		startEdge.Options = append(startEdge.Options, "Unique()")
		startNode.Edges = append(startNode.Edges, *startEdge)
	case "M2M":
		startNode.Edges = append(startNode.Edges, *startEdge)
		endNode.Edges = append(startNode.Edges, *endEdge)
	}
}

func (p *Parser) parseColumnType(column *types.Column) {
	column.DataType = p.h.Clean(column.DataType)

	if p.h.Contains(p.h.Pascal(column.DataType), "Enum") {
		p.cField.Type = "Enum"
		p.cField.EnumValues = p.h.CleanSplit(column.DataType, ",", "Enum", "(", ")")

		if p.config.Ent.Graphql != nil {
			values := p.h.MultiplyArray(p.cField.EnumValues, p.h.Pascal, p.h.UpperSnake)
			p.cField.Options = append(p.cField.Options, "NamedValues(\""+p.h.Join(values, "\", \"")+"\")")
		} else {
			values := p.h.MultiplyArray(p.cField.EnumValues, p.h.Pascal)
			p.cField.Options = append(p.cField.Options, "Values(\""+p.h.Join(values, "\", \"")+"\")")
		}
	} else {
		t := types.VuerdTypes[p.h.Lower(column.DataType)]

		if t != "" {
			p.cField.Type = EntTypesMap[t]
		} else {
			utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, type `%s` supported", p.cNode.Name, column.Name, column.DataType))
			p.cField.Type = column.DataType
		}
	}
}

func (p *Parser) parseColumnDefault(value string, option string) string {
	if p.cField.Type == "Enum" {
		if !p.h.InArray(p.cField.EnumValues, value) {
			utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, default value `%s` is not in enum", p.cNode.Name, p.cField.Name, value))
		} else if p.config.Ent.Graphql != nil {
			return option + "(\"" + p.h.UpperSnake(value) + "\")"
		} else {
			return option + "(\"" + p.h.Pascal(value) + "\")"
		}
	} else if p.cField.Type == "String" {
		return option + `("` + value + `")`
	} else {
		return option + `(` + value + `)`
	}
	return ""
}

func (p *Parser) appendFieldSkips(skips []string) {
	for _, skip := range skips {
		skip = SkipMap[skip]
		if !p.h.InArray(p.cField.Skips, skip) && skip != "" {
			p.cField.Skips = append(p.cField.Skips, skip)
		}
	}
}

func (p *Parser) appendFieldOptions(options []string, conditions ...bool) {
	for _, option := range options {
		p.appendFieldOption(option, conditions...)
	}
}

func (p *Parser) appendFieldOption(option string, conditions ...bool) {
	if len(conditions) > 0 {
		for _, condition := range conditions {
			if !condition {
				return
			}
		}
	}

	if !p.h.InArray(p.cField.Options, option) {
		p.cField.Options = append(p.cField.Options, option)
	}
}

func (p *Parser) parseColumnOptions(column *types.Column) {
	p.appendFieldOption("Unique()", column.Option.Unique)
	p.appendFieldOption("AutoIncrement()", column.Option.AutoIncrement)
	p.appendFieldOptions([]string{"Optional()", "Nillable()"}, !column.Option.NotNull)
	p.appendFieldOption(p.parseColumnDefault(p.h.Clean(column.Default), "Default"), len(p.h.Clean(column.Default)) > 0)
	p.appendFieldOptions(OptionMapDefault[p.cField.Name])
	p.addImports(ImportsDefault[p.cField.Name])

	if p.config.Ent.Graphql != nil {
		p.appendFieldSkips(SkipMapDefault[p.cField.Name])
		if p.h.InArray(ComparableTypesMap, p.cField.Type) && p.cField.Name != "password" {
			p.cField.Annotations = append(p.cField.Annotations, "entgql.OrderField(\""+p.h.UpperSnake(p.cField.Name)+"\")")
		}
	}

	cms := p.h.Split(column.Comment, "|")

	for _, cm := range cms {
		if p.h.Contains(cm, "=") {
			options := strings.Split(cm, "=")
			if options[1] == "upd" {
				p.appendFieldOption(p.parseColumnDefault(options[1], "UpdateDefault"), len(options[1]) > 0)
			} else if options[1] == "match" {
				p.appendFieldOption("Match(\""+options[1]+"\")", len(options[1]) > 0)
			} else if opts := RegexMap[options[0]]; len(opts) > 0 {
				for _, opt := range opts {
					if len(opt.Types) == 0 || p.h.InArray(opt.Types, p.cField.Type) {
						matched, err := regexp.MatchString(opt.Match, cm)
						utils.CatchError(utils.Fatal, err)
						if matched {
							reg, err := regexp.Compile(opt.Extract)
							utils.CatchError(utils.Fatal, err)
							values := reg.FindAllString(cm, -1)
							if options[0] == "skip" {
								p.appendFieldSkips(values)
							} else {
								p.appendFieldOption(opt.Option + p.h.Join(values, "(", ", ", ")"))
							}

						} else {
							utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, error in option `%s`", p.cNode.Name, p.cField.Name, cm))
						}
					}
				}
			} else {
				utils.CatchError(utils.Warninig, fmt.Errorf("table `%s`, column `%s`, unknown option `%s`", p.cNode.Name, p.cField.Name, cm))
			}
		}
	}
}

func (Parser) isMixin(name string) bool {
	return strings.Contains(strings.ToLower(name), "mixin")
}

func (p *Parser) table(id string) *types.Table {

	for i := range p.vuerd.TableState.Tables {
		if p.vuerd.TableState.Tables[i].Id == id {
			return &p.vuerd.TableState.Tables[i]
		}
	}

	return nil
}

func (p *Parser) node(id string) *Node {
	for _, n := range p.state.Nodes {
		if n.ID == id {
			return n
		}
	}

	for _, n := range p.state.Mixins {
		if n.ID == id {
			return n
		}
	}
	return nil
}

func (Parser) column(table *types.Table, id string) *types.Column {
	for i := 0; i < len(table.Columns); i++ {
		if table.Columns[i].Id == id {
			return &table.Columns[i]
		}
	}
	return nil
}

func (p *Parser) addImports(s string, nodes ...*Node) {
	if len(nodes) > 1 {
		for i := 0; i < len(nodes); i++ {
			if !p.h.InArray(nodes[i].Imports, s) {
				nodes[i].Imports = append(nodes[i].Imports, s)
			}
		}
	} else {
		if !p.h.InArray(p.cNode.Imports, s) {
			p.cNode.Imports = append(p.cNode.Imports, s)
		}
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

func (p *Parser) commentOptionValues(comment string, option string) []string {
	options := p.h.Split(comment, "|")

	if len(options) > 0 {
		for _, op := range options {
			if p.h.HasPreffix(op, option) {
				opts := p.h.Split(op, "=")
				return p.h.Split(opts[1], ",")
			}
		}
	}

	return []string{}
}

func (p *Parser) nodeExist(name string) bool {
	for i := range p.vuerd.TableState.Tables {
		if p.vuerd.TableState.Tables[i].Name == name {
			return true
		}
	}
	return false
}

func (p *Parser) validateM2MTable(table *types.Table) {
	found := []bool{false, false}
	for _, column := range table.Columns {
		if column.Ui.Fk {
			dr := p.commentOptionValues(column.Comment, "dr=")
			if len(dr) == 2 {
				if dr[0] == "To" {
					found[0] = true
				}
				if dr[0] == "From" {
					found[1] = true
				}
				if !p.nodeExist(dr[1]) {
					utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`, column `%s`, dr(%s,%s) %s does not exist", p.cNode.Name, column.Name, dr[0], dr[1], dr[1]))
				}
			}
		}
	}

	if !found[0] || !found[1] {
		utils.CatchError(utils.Fatal, fmt.Errorf("table `%s`,fk columns are not set correctly dr=(To/From, Type)", p.cNode.Name))
	}
}
