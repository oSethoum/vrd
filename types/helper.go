package types

import (
	"strconv"
	"strings"

	"github.com/codemodus/kace"
	"github.com/gertd/go-pluralize"
)

type Helper struct{}

func (Helper) Camel(s string) string {
	return kace.Camel(pluralize.NewClient().Singular(s))
}

func (Helper) Pascal(s string) string {
	return kace.Pascal(pluralize.NewClient().Singular(s))
}

func (Helper) Kebab(s string) string {
	return kace.Kebab(pluralize.NewClient().Singular(s))
}

func (Helper) Snake(s string) string {
	return kace.Snake(pluralize.NewClient().Singular(s))
}

func (Helper) Camels(s string) string {
	return kace.Camel(pluralize.NewClient().Plural(s))
}

func (h *Helper) MCamels(s string) string {
	return kace.Camel(h.MultiPlural(s))
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

func (Helper) MultiPlural(s string) string {
	s = kace.Snake(s)
	ss := strings.Split(s, "_")
	vv := ""
	for _, v := range ss {
		vv += kace.Pascal(pluralize.NewClient().Plural(v))
	}
	return vv
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

func (Helper) Lower(s string) string {
	return strings.ToLower(s)
}

func (Helper) Upper(s string) string {
	return strings.ToUpper(s)
}

func (Helper) UpperSnake(s string) string {
	return strings.ToUpper(kace.Snake(s))
}

func (Helper) Contains(s string, sub string) bool {
	return strings.Contains(s, sub)
}

func (Helper) Clean(target string, cs ...string) string {
	target = strings.ReplaceAll(target, " ", "")
	for _, c := range cs {
		target = strings.ReplaceAll(target, c, "")
	}
	return target
}

func (h Helper) Split(target string, ss string) []string {
	result := []string{}
	target = h.Clean(target)
	for _, s := range strings.Split(target, ss) {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func (h Helper) CleanSplit(target string, ss string, cs ...string) []string {
	return h.Split(h.Clean(target, cs...), ss)
}

func (h Helper) MultiplyArray(ss []string, fcts ...func(string) string) []string {
	result := []string{}
	for i := 0; i < len(ss); i++ {
		result = append(result, h.Multiply(ss[i], fcts...)...)
	}
	return result
}

func (h Helper) Multiply(s string, fcts ...func(string) string) []string {
	result := []string{}
	for i := 0; i < len(fcts); i++ {
		result = append(result, fcts[i](s))
	}
	return result
}

func (h Helper) HasPreffix(s string, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func (h Helper) HasSuffix(s string, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func (Helper) InArray(array []string, element string) bool {
	for _, v := range array {
		if element == v {
			return true
		}
	}
	return false
}
func (h Helper) ValueOfType(v string, t string) (interface{}, bool) {

	if h.InArray([]string{"String", "Byte", "Uint"}, t) {
		vv, err := strconv.ParseInt(v, 10, 32)
		return vv, err == nil
	} else if h.InArray([]string{"Float"}, t) {
		vv, err := strconv.ParseFloat(v, 32)
		return vv, err == nil
	} else if h.InArray([]string{"Bool"}, t) {
		vv, err := strconv.ParseBool(v)
		return vv, err == nil
	}

	return nil, false
}

func (Helper) StringNumberCompare(s1, s2 string) int {
	v1, _ := strconv.ParseFloat(s1, 32)
	v2, _ := strconv.ParseFloat(s2, 32)

	if v1 > v2 {
		return 1
	} else if v1 < v2 {
		return -1
	} else {
		return 0
	}
}

func (h *Helper) CleanArray(array []string) []string {
	result := []string{}
	for _, v := range array {
		if !h.InArray(result, v) {
			result = append(result, v)
		}
	}
	return result
}
