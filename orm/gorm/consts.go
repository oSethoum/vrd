package gorm

var GormTypesMap = map[string]string{
	"int":      "int",
	"long":     "int64",
	"float":    "float",
	"double":   "float64",
	"decimal":  "int",
	"boolean":  "bool",
	"string":   "string",
	"lob":      "string",
	"date":     "time.Time",
	"json":     "datatypes.JSON",
	"datetime": "time.Time",
	"time":     "time.Time",
}
