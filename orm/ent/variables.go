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
	"all":    "entgql.SkipAll",
	"type":   "entgql.SkipType",
	"where":  "entgql.SkipWhereInput",
	"create": "entgql.SkipMutationCreateInput",
	"update": "entgql.SkipMutationUpdateInput",
}

var RegexMap = map[string][]CommentOption{
	"skip": {
		{
			Match:   `^skip=(create|update|all|where|type)(,(create|update|all|where|type))*$`,
			Extract: `create|update|all|where|type`,
		},
	},

	"minLen": {
		{
			Types:   []string{"String", "Bytes"},
			Match:   `^minLen=(0|[1-9][0-9]*)$`,
			Extract: `[0-9]+`,
			Option:  "MinLen",
		},
	},

	"maxLen": {
		{
			Types:   []string{"String", "Bytes"},
			Match:   `^maxLen=(0|[1-9][0-9]*)$`,
			Extract: `[0-9]+`,
			Option:  "MaxLen",
		},
	},
	"min": {
		{
			Types:   []string{"Float"},
			Match:   `^min=(0|-?[1-9][0-9]*)(\.[0-9]+)?$`,
			Extract: `-?[0-9]+(\.[0-9]+)?`,
			Option:  "Min",
		},
		{
			Types:   []string{"Int"},
			Match:   `^min=(0|-?[1-9][0-9]*)$`,
			Extract: `-?[0-9]+`,
			Option:  "Min",
		},
		{
			Types:   []string{"Uint"},
			Match:   `^min=(0|[1-9][0-9]*)$`,
			Extract: `[0-9]+`,
			Option:  "Min",
		},
	},
	"max": {
		{
			Types:   []string{"Float"},
			Match:   `^max=(0|-?[1-9][0-9]*)(\.[0-9]+)?$`,
			Extract: `-?[0-9]+(\.[0-9]+)?`,
			Option:  "Max",
		},
		{
			Types:   []string{"Int"},
			Match:   `^max=(0|-?[1-9][0-9]*)$`,
			Extract: `-?[0-9]+`,
			Option:  "Max",
		},
		{
			Types:   []string{"Uint"},
			Match:   `^max=(0|[1-9][0-9]*)$`,
			Extract: `[0-9]+`,
			Option:  "Max",
		},
	},
	"range": {
		{
			Types:   []string{"Float"},
			Match:   `^range=(0|-?[1-9][0-9]*)(\.[0-9]+)?,(0|-?[1-9][0-9]*)(\.[0-9]+)?$`,
			Extract: `-?[0-9]+(\.[0-9]+)?`,
			Option:  "Range",
		},
		{
			Types:   []string{"Int"},
			Match:   `^range=(0|-?[1-9][0-9]*),(0|-?[1-9][0-9]*)$`,
			Extract: `-?[0-9]+`,
		},
		{
			Types:   []string{"Uint"},
			Match:   `^range=(0|[1-9][0-9]*),(0|[1-9][0-9]*)$`,
			Extract: `[0-9]+`,
			Option:  "Range",
		},
	},
}

var SkipMapDefault = map[string][]string{
	"created_at": {"create", "update"},
	"updated_at": {"create", "update"},
	"password":   {"where"},
}

var OptionMapDefault = map[string][]string{
	"password":   {"Sensitive()"},
	"created_at": {"Default(\"time.Now\")"},
	"updated_at": {"Default(\"time.Now\")", "UpdateDefault(\"time.Now\")"},
}
