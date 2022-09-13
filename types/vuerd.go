package types

import (
	"strings"

	"github.com/codemodus/kace"
	"github.com/gertd/go-pluralize"
)

type State struct {
	TableState        TableState        `json:"table"`
	RelationshipState RelationshipState `json:"relationShip"`
}

type TableState struct {
	Tables  []Table `json:"tables"`
	Indexes []Index `json:"indexes"`
}

type Table struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Comment string   `json:"comment"`
	Columns []Column `json:"columns"`
	Helper
}

type Column struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Comment  string       `json:"comment"`
	DataType string       `json:"dataType"`
	Default  string       `json:"default"`
	Option   ColumnOption `json:"option"`
	Ui       ColumnUi     `json:"ui"`
}

type ColumnOption struct {
	AutoIncrement bool `json:"autoIncrement"`
	PrimaryKey    bool `json:"primaryKey"`
	Unique        bool `json:"unique"`
	NotNull       bool `json:"notNull"`
}

type ColumnUi struct {
	Pk  bool `json:"pk"`
	Fk  bool `json:"fk"`
	Pfk bool `json:"pfk"`
}

type RelationshipState struct {
	Relationships []Relationship `json:"relationships"`
}

type Relationship struct {
	Id               string            `json:"id"`
	Identification   bool              `json:"identification"`
	RelationshipType string            `json:"relationshipType"` // RelationshipType = ZeroN | OneN | ZeroOne | OneOnly
	Start            RelationshipPoint `json:"start"`
	End              RelationshipPoint `json:"end"`
	ConstraintName   string            `json:"constraintName?"`
}

type RelationshipPoint struct {
	TableId   string   `json:"tableId"`
	ColumnIds []string `json:"columnIds"`
}

type MemoState struct {
	Memos []Memo `json:"memos"`
}

type Memo struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type Index struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	TableId string        `json:"tableId"`
	Columns []IndexColumn `json:"columns"`
	Unique  bool          `json:"unique"`
}

// OrderType = ASC | DESC
type IndexColumn struct {
	Id        string `json:"id"`
	OrderType string `json:"orderType"`
}

type File struct {
	Path   string `json:"path"`
	Buffer string `json:"buffer"`
}

var VuerdTypes = map[string]string{
	"bfile":              "lob",
	"bigint":             "long",
	"bigserial":          "long",
	"binary":             "string",
	"binary_double":      "double",
	"binary_float":       "float",
	"bit":                "int",
	"blob":               "lob",
	"bool":               "boolean",
	"boolean":            "boolean",
	"box":                "string",
	"bytea":              "string",
	"char":               "string",
	"character":          "string",
	"cidr":               "string",
	"circle":             "string",
	"clob":               "lob",
	"date":               "date",
	"datetime":           "datetime",
	"datetime2":          "datetime",
	"datetimeoffset":     "datetime",
	"dec":                "decimal",
	"decimal":            "decimal",
	"double":             "double",
	"enum":               "string",
	"fixed":              "decimal",
	"float":              "float",
	"float4":             "float",
	"float8":             "double",
	"geography":          "string",
	"geometry":           "string",
	"geometrycollection": "string",
	"image":              "lob",
	"inet":               "string",
	"int":                "int",
	"int2":               "int",
	"int4":               "int",
	"int8":               "long",
	"integer":            "int",
	"interval":           "time",
	"json":               "json",
	"jsonb":              "lob",
	"line":               "string",
	"linestring":         "string",
	"long":               "lob",
	"longblob":           "lob",
	"longtext":           "lob",
	"lseg":               "string",
	"macaddr":            "string",
	"macaddr8":           "string",
	"mediumblob":         "lob",
	"mediumint":          "int",
	"mediumtext":         "lob",
	"money":              "double",
	"multilinestring":    "string",
	"multipoint":         "string",
	"multipolygon":       "string",
	"nchar":              "string",
	"nclob":              "lob",
	"ntext":              "lob",
	"number":             "long",
	"numeric":            "decimal",
	"nvarchar":           "string",
	"nvarchar2":          "string",
	"path":               "string",
	"pg_lsn":             "int",
	"point":              "string",
	"polygon":            "string",
	"raw":                "lob",
	"real":               "double",
	"serial":             "int",
	"serial2":            "int",
	"serial4":            "int",
	"serial8":            "long",
	"set":                "string",
	"smalldatetime":      "datetime",
	"smallint":           "int",
	"smallmoney":         "float",
	"smallserial":        "int",
	"sql_variant":        "string",
	"text":               "lob",
	"time":               "time",
	"timestamp":          "datetime",
	"timestamptz":        "datetime",
	"timetz":             "time",
	"tinyblob":           "lob",
	"tinyint":            "int",
	"tinytext":           "lob",
	"tsquery":            "string",
	"tsvector":           "string",
	"txid_snapshot":      "string",
	"uniqueidentifier":   "string",
	"uritype":            "string",
	"uuid":               "uuid",
	"varbinary":          "string",
	"varbit":             "int",
	"varchar":            "string",
	"varchar2":           "string",
	"xml":                "lob",
	"xmltype":            "string",
	"year":               "int",
}

type Helper struct{}

func (Helper) Camel(s string) string {
	return kace.Camel(s)
}

func (Helper) Pascal(s string) string {
	return kace.Pascal(s)
}

func (Helper) Kebab(s string) string {
	return kace.Kebab(s)
}

func (Helper) Snake(s string) string {
	return kace.Snake(s)
}

func (Helper) Camels(s string) string {
	return kace.Camel(pluralize.NewClient().Plural(s))
}

func (Helper) Pascals(s string) string {
	return kace.Pascal(pluralize.NewClient().Plural(s))
}

func (Helper) Kebabs(s string) string {
	return kace.Kebab(pluralize.NewClient().Plural(s))
}

func (Helper) Snakes(s string) string {
	return kace.Snake(pluralize.NewClient().Plural(s))
}

func (Helper) Plural(s string) string {
	return pluralize.NewClient().Plural(s)
}

func (Helper) Singular(s string) string {
	return pluralize.NewClient().Singular(s)
}

func (Helper) Join(ss []string, args ...string) string {
	if len(ss) == 0 {
		return ""
	}
	switch len(args) {
	case 1:
		return strings.Join(ss, args[0])
	case 2:
		return args[0] + strings.Join(ss, args[1])
	default:
		return args[0] + strings.Join(ss, args[1]) + args[2]
	}
}
