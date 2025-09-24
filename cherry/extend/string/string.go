package cherryString

import (
	"encoding/json"
	"strconv"
)

func ToString(value any) string {
	ret := ""

	if value == nil {
		return ret
	}

	switch t := value.(type) {
	case string:
		ret = t
	case int:
		ret = strconv.Itoa(t)
	case int32:
		ret = strconv.Itoa(int(t))
	case int64:
		ret = strconv.FormatInt(t, 10)
	case uint:
		ret = strconv.Itoa(int(t))
	case uint32:
		ret = strconv.Itoa(int(t))
	case uint64:
		ret = strconv.Itoa(int(t))
	default:
		v, _ := json.Marshal(t)
		ret = string(v)
	}

	return ret
}

func ToStringSlice(val []any) []string {
	var result []string

	for _, item := range val {
		v, ok := item.(string)
		if ok {
			result = append(result, v)
		}
	}

	return result
}
