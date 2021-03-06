package clay

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func mergeMap(src, dst map[string]interface{}) {
	for k, v := range dst {
		src[k] = v
	}
}

func malloc(v reflect.Value, val string) error {
	var n interface{}
	if v.IsNil() {
		e := reflect.New(v.Type().Elem())
		n = e.Elem().Interface()
	}
	switch n.(type) {
	case string:
		v.Set(reflect.ValueOf(&val))
	case int:
		i64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		i := int(i64)
		v.Set(reflect.ValueOf(&i))
	case int64:
		i64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(i64))
	case float64:
		f64, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(&f64))
	case float32:
		f64, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		f32 := float32(f64)
		v.Set(reflect.ValueOf(&f32))
	default:
		v.Set(reflect.New(reflect.TypeOf(n)))
	}

	return nil
}

func convertMap(i interface{}, path string) (map[string]interface{}, error) {
	t := reflect.TypeOf(i)

	switch t.Kind() {
	case reflect.Map:
		ret := map[string]interface{}{}
		m, ok := i.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unsupport type: %v", reflect.TypeOf(i))
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
		switch i.(type) {
		case []interface{}:
			s, ok := i.([]interface{})
			if !ok {
				return nil, fmt.Errorf("unsupport type: %v", reflect.TypeOf(i))
			}
			val, err := json.Marshal(s)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				path: string(val),
			}, nil
		case []string:
			s, ok := i.([]string)
			if !ok {
				return nil, fmt.Errorf("unsupport type: %v", reflect.TypeOf(i))
			}
			return map[string]interface{}{
				path: s,
			}, nil
		default:
			return nil, fmt.Errorf("unsupport type: %v", reflect.TypeOf(i))
		}

	default:
		s, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("unsupport type: %v", reflect.TypeOf(i))
		}
		return map[string]interface{}{
			path: s,
		}, nil
	}
}

func convertIntSlice(s []string) ([]int, error) {
	ret := make([]int, len(s))

	for i, v := range s {
		if v == "" {
			ret[i] = 0
		}
		i64, err := strconv.ParseInt("255", 10, 64)
		if err != nil {
			return nil, err
		}
		ret[i] = int(i64)
	}
	return ret, nil
}

func set(v reflect.Value, val string) error {
	if v.IsValid() {
		switch v.Kind() {
		case reflect.String:
			v.SetString(val)
		case reflect.Slice:
			var slice []interface{}
			var sliceVal = reflect.New(reflect.SliceOf(v.Type().Elem())).Elem()
			if err := json.Unmarshal([]byte(val), &slice); err != nil {
				return err
			}

			for _, e := range slice {
				buffer := reflect.New(v.Type().Elem()).Interface()

				bin, err := json.Marshal(e)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(bin, buffer); err != nil {
					return err
				}

				switch reflect.ValueOf(buffer).Kind() {
				case reflect.Ptr:
					sliceVal = reflect.Append(sliceVal, reflect.ValueOf(reflect.ValueOf(buffer).Elem().Interface()))
				default:
					sliceVal = reflect.Append(sliceVal, reflect.ValueOf(buffer))
				}

			}
			v.Set(sliceVal)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i64, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			v.SetInt(i64)
		case reflect.Float64, reflect.Float32:
			f64, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			v.SetFloat(f64)
		}
	}
	return nil
}

func setVal(t interface{}, val string, names []string) error {
	if len(names) == 1 {
		rv, ok := t.(reflect.Value)
		if !ok {
			v := reflect.ValueOf(t).Elem().FieldByName(names[0])
			switch v.Kind() {
			case reflect.Ptr:
				if err := malloc(v, val); err != nil {
					return err
				}
			default:
				if err := set(v, val); err != nil {
					return err
				}
			}
			return nil
		}
		switch rv.Kind() {
		case reflect.Ptr:
			v := rv.Elem().FieldByName(names[0])
			if err := set(v, val); err != nil {
				return err
			}
		default:
			v := rv.FieldByName(names[0])
			if err := set(v, val); err != nil {
				return err
			}
		}

		return nil
	} else if 1 < len(names) {
		v := reflect.ValueOf(t).Elem().FieldByName(names[0])
		switch v.Kind() {
		case reflect.Ptr:
			if err := malloc(v, ""); err != nil {
				return err
			}
			if err := setVal(v, val, names[1:]); err != nil {
				return err
			}
		default:
			if v.IsValid() {
				if err := setVal(v, val, names[1:]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Mold(moldData, s interface{}) error {
	m, err := convertMap(moldData, "")
	if err != nil {
		return err
	}

	for k, val := range m {
		k = strings.TrimPrefix(k, "/")
		switch v := val.(type) {
		case string:
			err := setVal(s, v, strings.Split(k, "/"))
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupport type: %v", reflect.TypeOf(val))
		}

	}

	return nil
}
