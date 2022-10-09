package ent

import (
	"fmt"
	"log"
	"strings"

	"vrd/config"
	"vrd/types"

	"github.com/codemodus/kace"
	"github.com/gertd/go-pluralize"
)

//TODO: support onupdate and ondelete for sql

func Parse(state types.State, config config.Config) State {

	st := State{
		Nodes:  []Node{},
		Mixins: make(map[string]Mixin),
	}

	if len(state.TableState.Tables) == 0 {
		log.Fatal("Parser: no table found")
	}

	// parsing mixins
	for _, table := range state.TableState.Tables {
		if sClean(table.Name) == "" {
			log.Fatalf("Parser: Error Table without name id = {%s}", table.Id)
		}
		if strings.Contains(table.Name, "Mixin") {
			mixin := parseTableMixin(&state, &table, &config)
			st.Mixins[mixin.Alias] = *mixin
		}
	}

	// parsing tables
	for _, table := range state.TableState.Tables {
		if sClean(table.Name) == "" {
			log.Fatalf("Parser: Error Table without name id = {%s}", table.Id)
		}

		if !strings.Contains(table.Name, "Mixin") {
			node := parseTableNode(&state, st.Mixins, &table, &config)
			if node != nil {
				st.Nodes = append(st.Nodes, *node)
			}
		}
	}
	return st
}

func parseTableNode(state *types.State, mixins map[string]Mixin, table *types.Table, config *config.Config) *Node {

	node := Node{
		ID:      table.Id,
		Name:    table.Name,
		Comment: sClean(table.Comment),
		Imports: []string{
			"\t\"entgo.io/ent\"",
			"\t\"entgo.io/ent/dialect/entsql\"",
			"\t\"entgo.io/ent/schema\"",
		},
	}

	if config.Ent.Privacy && config.Ent.PrivacyNode {
		node.Imports = append(node.Imports, "\t\""+config.Ent.Package+"/ent/privacy\"")
		node.Imports = append(node.Imports, "\t\""+config.Ent.Package+"/auth\"")
	}

	var isM2M bool

	tableComment := sClean(table.Comment)
	coptions := strings.Split(tableComment, "|")
	for _, cop := range coptions {
		if cop == "" {
			continue
		}
		if strings.ToLower(cop) == "m2m" {
			// handle the case of m2m
			isM2M = true
			if len(table.Columns) == 2 {
				return nil
			}
			continue
		}
		cops := strings.Split(cop, "=")
		switch cops[0] {
		case "nm":
			{
				node.TableName = cops[1]
				break
			}

		case "mxs":
			{

				mxs := strings.Split(strings.ReplaceAll(strings.ReplaceAll(cops[1], "(", ""), ")", ""), ",")
				for _, mx := range mxs {
					node.Mixins = append(node.Mixins, mixins[mx].Name)
				}
				break
			}

		default:
			fmt.Printf("Parser: warning unknwon param in table's comment: %s, param: %s", table.Name, cops[0])
		}
	}

	if config.Ent.Graphql {
		node.Imports = append(node.Imports, "\t\"entgo.io/contrib/entgql\"")
	}

	node.Fields = make([]Field, 0)

	for _, column := range table.Columns {

		columnErrors(column, table.Name)

		if isM2M && column.Ui.Fk {
			field := parseColumn(table, &column, config)
			node.Fields = append(node.Fields, field)
		}

		if !column.Ui.Pk && !column.Ui.Fk && !column.Ui.Pfk {
			field := parseColumn(table, &column, config)

			if strings.Contains(strings.ToLower(field.Type), "time") && !in("\t\"time\"", node.Imports) {
				node.Imports = append(node.Imports, "\t\"time\"")
			}
			node.Fields = append(node.Fields, field)
		}

	}

	node.Edges = make([]Edge, 0)

	for _, relationship := range state.RelationshipState.Relationships {
		edge := Edge{
			ID: relationship.Id,
		}
		options := []string{}

		if relationship.Start.TableId == table.Id {
			endColumn := findColumn(state, relationship.End.TableId, relationship.End.ColumnIds[0])
			endTable := findTable(state, relationship.End.TableId)
			edge.Direction = "To"

			if strings.Contains(strings.ToLower(endTable.Comment), "m2m") {
				if sClean(endTable.Comment) == "" {
					log.Fatalf("Parsing: Error M2M Table %s, Column %s cannot be empty", endTable.Name, endColumn.Name)
				}
				cms := strings.Split(sClean(endColumn.Comment), ",")
				edge.Direction = cms[0]
				edge.Node = cms[2]
				edge.Name = cms[1]
				tb := strings.Split(strings.Split(endTable.Comment, "|")[1], "=")[1]
				if tb == "" {
					tb = multiPlural(endTable.Name)
				}

				if len(endTable.Columns) > 2 {
					edge.Options = append(edge.Options, fmt.Sprintf("Through(\"%s\", %s.Type)", tb, endTable.Name))
				}

				if edge.Direction == "From" {
					edge.Options = append(edge.Options, fmt.Sprintf("Ref(\"%s\")", cms[3]))
				}
			} else {
				edge.Node = endTable.Name
			}

			if len(strings.Split(endColumn.Comment, "|")) < 2 {
				switch relationship.RelationshipType {
				case "ZeroN", "OneN":
					edge.Name = kace.Camel(pluralize.NewClient().Plural(edge.Node))
					edge.Annotations = append(edge.Annotations, "entgql.RelayConnection()")
				case "ZeroOne", "OneOnly":
					edge.Name = kace.Camel(edge.Node)
				}
			} else {
				edge.Name = sClean(strings.Split(endColumn.Comment, "|")[1])
			}

			node.Edges = append(node.Edges, edge)
		}

		if relationship.End.TableId == table.Id {
			endColumn := findColumn(state, relationship.End.TableId, relationship.End.ColumnIds[0])
			if isM2M {
				edge.Direction = "To"
			} else {
				edge.Direction = "From"
			}
			edge.Node = findTable(state, relationship.Start.TableId).Name

			if len(strings.Split(endColumn.Comment, "|")) < 2 {
				etName := findTable(state, relationship.End.TableId).Name
				edge.Name = kace.Camel(edge.Node)

				switch relationship.RelationshipType {
				case "OneN", "ZeroN":
					edge.Reference = kace.Camel(pluralize.NewClient().Plural(etName))
				case "OneOnly", "ZeroOne":
					edge.Reference = kace.Camel(etName)
				}
			} else {
				edge.Name = sClean(strings.Split(endColumn.Comment, "|")[0])
				edge.Reference = sClean(strings.Split(endColumn.Comment, "|")[1])
			}

			if isM2M {
				options = append(options, "Unique()", "Required()", fmt.Sprintf("Field(\"%s\")", endColumn.Name))

			} else {
				switch relationship.RelationshipType {
				case "ZeroOne", "ZeroN":
					options = append(options, "Unique()")
				case "OneOnly", "OneN":
					options = append(options, "Unique()", "Required()")
				}
			}

			if !isM2M {
				options = append(options, "Ref(\""+edge.Reference+"\")")
			}

			edge.Options = options
			node.Edges = append(node.Edges, edge)
		}

	}

	if len(node.Edges) > 0 {
		node.Imports = append(node.Imports, "\t\"entgo.io/ent/schema/edge\"")
	}

	if len(node.Fields) > 0 {
		node.Imports = append(node.Imports, "\t\"entgo.io/ent/schema/field\"")
	}

	if node.TableName == "" {
		node.TableName = kace.Snake(node.MultiPlural(node.Name))
	}

	node.Annotations = []string{
		"\t\tentsql.Annotation{Table: \"" + node.TableName + "\"}",
	}

	if config.Ent.Graphql {
		node.Annotations = append(node.Annotations, []string{
			"\t\tentgql.QueryField(\"" + kace.Camel(node.TableName) + "\")",
			"\t\tentgql.RelayConnection()",
			"\t\tentgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate())",
		}...)
	}

	return &node
}

func parseTableMixin(state *types.State, table *types.Table, config *config.Config) *Mixin {
	mixin := Mixin{
		ID:      table.Id,
		Name:    table.Name,
		Comment: sClean(table.Comment),
		Imports: []string{
			"\t\"entgo.io/ent\"",
			"\t\"entgo.io/ent/schema/mixin\"",
		},
	}

	if config.Ent.Graphql {
		mixin.Imports = append(mixin.Imports, "\t\"entgo.io/contrib/entgql\"")
	}

	var isM2M bool

	coptions := strings.Split(sClean(table.Comment), "|")
	for _, cop := range coptions {
		if strings.ToLower(cop) == "m2m" {
			// handle the case of m2m
			isM2M = true
			if len(table.Columns) == 2 {
				return nil
			}
			continue
		}
		cops := strings.Split(cop, "=")
		switch cops[0] {
		case "nx":
			{
				mixin.Alias = cops[1]
				break
			}
		default:
			fmt.Printf("Parser: warning unknwon param in table's comment: %s, param: %s", table.Name, cops[0])
		}
	}

	mixin.Fields = make([]Field, 0)

	for _, column := range table.Columns {

		columnErrors(column, table.Name)

		if isM2M && column.Ui.Fk {
			field := parseColumn(table, &column, config)
			mixin.Fields = append(mixin.Fields, field)
		}

		if !column.Ui.Pk && !column.Ui.Fk && !column.Ui.Pfk {
			field := parseColumn(table, &column, config)

			if strings.Contains(strings.ToLower(field.Type), "time") && !in("\t\"time\"", mixin.Imports) {
				mixin.Imports = append(mixin.Imports, "\t\"time\"")
			}
			mixin.Fields = append(mixin.Fields, field)
		}

	}

	mixin.Edges = make([]Edge, 0)

	for _, relationship := range state.RelationshipState.Relationships {
		edge := Edge{
			ID: relationship.Id,
		}
		options := []string{}

		if relationship.Start.TableId == table.Id {
			endColumn := findColumn(state, relationship.End.TableId, relationship.End.ColumnIds[0])
			endTable := findTable(state, relationship.End.TableId)
			edge.Direction = "To"

			if strings.Contains(strings.ToLower(endTable.Comment), "m2m") {
				if sClean(endTable.Comment) == "" {
					log.Fatalf("Parsing: Error M2M Table %s, Column %s cannot be empty", endTable.Name, endColumn.Name)
				}
				cms := strings.Split(sClean(endColumn.Comment), ",")
				edge.Direction = cms[0]
				edge.Node = cms[2]
				edge.Name = cms[1]
				tb := strings.Split(strings.Split(endTable.Comment, "|")[1], "=")[1]
				if tb == "" {
					tb = multiPlural(endTable.Name)
				}

				if len(endTable.Columns) > 2 {
					edge.Options = append(edge.Options, fmt.Sprintf("Through(\"%s\", %s.Type)", tb, endTable.Name))
				}

				if edge.Direction == "From" {
					edge.Options = append(edge.Options, fmt.Sprintf("Ref(\"%s\")", cms[3]))
				}
			} else {
				edge.Node = endTable.Name
			}

			if len(strings.Split(endColumn.Comment, "|")) < 2 {
				switch relationship.RelationshipType {
				case "ZeroN", "OneN":
					edge.Name = kace.Camel(pluralize.NewClient().Plural(edge.Node))
				case "ZeroOne", "OneOnly":
					edge.Name = kace.Camel(edge.Node)
				}
			} else {
				edge.Name = sClean(strings.Split(endColumn.Comment, "|")[1])
			}

			mixin.Edges = append(mixin.Edges, edge)
		}

		if relationship.End.TableId == table.Id {
			endColumn := findColumn(state, relationship.End.TableId, relationship.End.ColumnIds[0])
			edge.Direction = "From"
			edge.Node = findTable(state, relationship.Start.TableId).Name

			if len(strings.Split(endColumn.Comment, "|")) < 2 {
				etName := findTable(state, relationship.End.TableId).Name
				edge.Name = kace.Camel(edge.Node)

				switch relationship.RelationshipType {
				case "OneN", "ZeroN":
					edge.Reference = kace.Camel(pluralize.NewClient().Plural(etName))
				case "OneOnly", "ZeroOne":
					edge.Reference = kace.Camel(etName)
				}
			} else {
				edge.Name = sClean(strings.Split(endColumn.Comment, "|")[0])
				edge.Reference = sClean(strings.Split(endColumn.Comment, "|")[1])
			}

			if isM2M {
				options = append(options, "Unique()", "Required()", fmt.Sprintf("Field(\"%s\")", endColumn.Name))

			} else {
				switch relationship.RelationshipType {
				case "ZeroOne", "ZeroN":
					options = append(options, "Unique()")
				case "OneOnly", "OneN":
					options = append(options, "Unique()", "Required()")
				}
			}

			if !isM2M {
				options = append(options, "Ref(\""+edge.Reference+"\")")
			}

			edge.Options = options
			mixin.Edges = append(mixin.Edges, edge)
		}

	}

	if len(mixin.Edges) > 0 {
		mixin.Imports = append(mixin.Imports, "\t\"entgo.io/ent/schema/edge\"")
	}

	if len(mixin.Fields) > 0 {
		mixin.Imports = append(mixin.Imports, "\t\"entgo.io/ent/schema/field\"")
	}

	return &mixin
}

func findColumn(state *types.State, tableId string, columnId string) *types.Column {
	for i, table := range state.TableState.Tables {
		for j, column := range table.Columns {
			if column.Id == columnId {
				return &state.TableState.Tables[i].Columns[j]
			}
		}
	}
	return nil
}

func findTable(state *types.State, tableId string) *types.Table {
	for i, table := range state.TableState.Tables {
		if table.Id == tableId {
			return &state.TableState.Tables[i]
		}
	}
	return nil
}

func parseColumn(table *types.Table, column *types.Column, config *config.Config) Field {
	datatype := parseType(column.DataType)
	options := []string{}
	annotations := []string{}
	gqlSkips := []string{}
	enumValues := []string{}
	defaultValue := parseDefault(datatype, column.Default)
	commentOptions := strings.Split(sClean(column.Comment), "|")

	for _, cop := range commentOptions {
		if strings.Contains(cop, "skip") {
			cop = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(cop, "skip", ""), "(", ""), ")", "")
			skips := strings.Split(cop, ",")
			for _, skip := range skips {
				switch skip {
				case "all":
					{
						gqlSkips = append(gqlSkips, "entgql.SkipAll")
						break
					}
				case "type":
					{
						gqlSkips = append(gqlSkips, "entgql.SkipType")
						break
					}

				case "create":
					{
						gqlSkips = append(gqlSkips, "entgql.SkipMutationCreateInput")
						break
					}
				case "update":
					{
						gqlSkips = append(gqlSkips, "entgql.SkipMutationUpdateInput")
						break
					}
				case "where":
					{
						gqlSkips = append(gqlSkips, "entgql.SkipWhereInput")
						break
					}
				default:
					fmt.Printf("Parser: warning unknwon param in table: %s, column: %s, param: %s", table.Name, column.Name, skip)
				}

			}
		}

		if strings.Contains(cop, "=") {
			arr := strings.Split(cop, "=")

			if arr[0] == "upd" {
				if datatype == "String" || datatype == "Bytes" || datatype == "Enum" {
					options = append(options, "UpdateDefault(\""+arr[1]+"\")")
				} else {
					options = append(options, "UpdateDefault("+arr[1]+")")
				}
			}

			if datatype == "String" || datatype == "Bytes" {
				switch arr[0] {
				case "min":
					{
						options = append(options, "MinLen("+arr[1]+")")
						break
					}

				case "max":
					{
						options = append(options, "MaxLen("+arr[1]+")")
						break
					}

				default:
					fmt.Printf("Parser: warning unknwon param in table: %s, column: %s, param: %s", table.Name, column.Name, arr[0])
				}

			}

			if datatype == "String" && arr[0] == "match" {
				options = append(options, "Match(\""+arr[1]+"\")")
			}

			if strings.Contains(datatype, "Int") || strings.Contains(datatype, "Float") {
				switch arr[0] {
				case "min":
					{
						options = append(options, "Min("+arr[1]+")")
						break
					}

				case "max":
					{
						options = append(options, "Max("+arr[1]+")")
						break
					}

				case "range":
					{

						mnx := strings.Split(arr[1], ",")
						options = append(options, "Range("+mnx[0]+", "+mnx[1]+")")
						break
					}
				default:
					fmt.Printf("Parser: warning unknwon param in table: %s, column: %s, param: %s", table.Name, column.Name, arr[0])
				}

			}
		}
	}

	if datatype == "String" || datatype == "Bytes" {
		if strings.Contains(column.Comment, "-nem") {
			options = append(options, "NotEmpty()")
		}
	}
	if strings.Contains(datatype, "Int") || strings.Contains(datatype, "Float") {
		if strings.Contains(column.Comment, "-pos") {
			options = append(options, "Positive()")
		}
		if strings.Contains(column.Comment, "-neg") {
			options = append(options, "Negative()")
		}
		if strings.Contains(column.Comment, "-nneg") {
			options = append(options, "NonNegative()")
		}
	}

	if column.Option.Unique {
		options = append(options, "Unique()")
	}

	if column.Option.AutoIncrement {
		options = append(options, "AutoIncrement()")
	}

	if datatype == "Enum" {
		enumValues = parseEnumValues(column.DataType)
		if config.Ent.Graphql {
			enOp := []string{}
			for _, v := range enumValues {
				enOp = append(enOp, kace.Pascal(strings.ToLower(v)), kace.SnakeUpper(v))
			}
			options = append(options, "NamedValues(\""+strings.Join(enOp, "\", \"")+"\")")
		} else {
			options = append(options, "Values(\""+kace.Pascal(strings.Join(enumValues, "\", \""))+"\")")
		}
	}

	if len(sClean(defaultValue)) > 0 {
		if datatype == "Enum" {
			if config.Ent.Graphql {
				options = append(options, "Default("+kace.SnakeUpper(defaultValue)+")")

			} else {
				options = append(options, "Default("+kace.Pascal(defaultValue)+")")
			}
		} else {
			options = append(options, "Default("+defaultValue+")")
		}
	}

	if strings.Contains(column.Comment, "-im") {
		options = append(options, "Immutable()")
	}

	if strings.Contains(column.Comment, "-s") {
		options = append(options, "Sensitive()")
		if config.Ent.Graphql {
			annotations = append(annotations, []string{
				"entgql.Skip(entgql.SkipType, entgql.SkipWhereInput)",
			}...)
		}
	} else {
		if config.Ent.Graphql && in(datatype, ComparableTypes) {
			annotations = append(annotations, []string{
				"entgql.OrderField(\"" + strings.ToUpper(column.Name) + "\")",
			}...)
		}
	}

	if !column.Option.NotNull {
		options = append(options, "Optional()", "Nillable()")
	}

	if strings.Contains(column.Comment, "-op") && column.Option.NotNull {
		options = append(options, "Optional()")
	}

	if len(gqlSkips) > 0 {
		annotations = append(annotations, "entgql.Skip("+strings.Join(gqlSkips, ", ")+")")
	}

	options = append(options, "Annotations("+strings.Join(annotations, ", ")+")")

	return Field{
		ID:          column.Id,
		Name:        column.Name,
		Type:        datatype,
		Options:     options,
		EnumValues:  enumValues,
		Default:     defaultValue,
		Comment:     column.Comment,
		Annotations: annotations,
	}
}

func parseType(datatype string) string {
	if strings.Contains(strings.ToUpper(datatype), "ENUM") {
		return "Enum"
	}
	t := EntTypes[types.VuerdTypes[strings.ToLower(datatype)]]

	if t != "" {
		return t
	}

	return datatype
}

func parseEnumValues(datatype string) []string {
	if !strings.Contains(datatype, "ENUM") {
		return nil
	}
	return strings.Split(sClean(strings.Split(strings.Split(datatype, "(")[1], ")")[0]), ",")
}

func parseDefault(datatype string, defaultVal string) string {
	if (datatype == "String" || datatype == "Enum") && len(sClean(defaultVal)) > 0 {
		return "\"" + defaultVal + "\""
	}
	return defaultVal
}

// Helper
func in[T int | float32 | string](value T, list []T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func columnErrors(column types.Column, table string) {
	if sClean(column.Name) == "" {
		log.Fatalf("Parser: Error column has no name %s", table)
	}

	if sClean(column.DataType) == "" {
		log.Fatalf("Parser: Error column has no datatype %s.%s", table, column.Name)
	}
}

func sClean(s string) string { return strings.ReplaceAll(s, " ", "") }

func multiPlural(s string) string {
	s = kace.Snake(s)
	ss := strings.Split(s, "_")
	vv := ""
	for _, v := range ss {
		vv += kace.Pascal(pluralize.NewClient().Plural(v))
	}
	return vv
}
