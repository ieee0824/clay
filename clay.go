package clay

import (
	"reflect"
	"errors"
	"fmt"
)

var ConvertMap = convertMap

func mergeMap(src, dst map[string]interface{}) {
	for k, v := range dst {
		src[k] = v
	}
}

func convertMap(i interface{}, path string) (map[string]interface{}, error) {
	t := reflect.TypeOf(i)

	switch t.Kind() {
	case reflect.Map:
		ret := map[string]interface{}{}
		m, ok := i.(map[string]interface{})
		if !ok {
			return nil, errors.New("")
		}

		for k, v := range m {
			result, err := convertMap(v, fmt.Sprintf("%s/%s", path, k))
			if err != nil {
				return nil, err
			}
			mergeMap(ret, result)
		}
		return ret, nil

	case reflect.Slice:
		s, ok := i.([]string)
		if !ok {
			return nil, errors.New("")
		}
		return map[string] interface{} {
			path: s,
		}, nil
	default:
		s, ok := i.(string)
		if !ok {
			return nil, errors.New("")
		}
		return map[string]interface{} {
			path: s,
		}, nil
	}
}