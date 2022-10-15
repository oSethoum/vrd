package ent

var EntTypesMap = map[string]string{
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

var ComparableTypesMap = []string{
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

var OptionsMap = map[string]string{
	"-nem":  "NotEmpty()",
	"-pos":  "Positive()",
	"-neg":  "Negative()",
	"-nneg": "NonNegative()",
	"-im":   "Immutable()",
	"-s":    "Sensitive()",
	"-op":   "Optional()",
}

var SkipMap = map[string]string{
	"all":    "entgql.SkipAll()",
	"type":   "entgql.SkipType()",
	"where":  "entgql.SkipWhereInput()",
	"create": "entgql.SkipMutationCreateInput()",
	"update": "entgql.SkipMutationUpdateInput()",
}
