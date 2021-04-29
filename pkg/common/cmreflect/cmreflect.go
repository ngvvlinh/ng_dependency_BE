package cmreflect

import (
	"fmt"
	"reflect"
	"strings"
)

func ReadTag(f reflect.StructField, tag string) (string, bool) {
	val, ok := f.Tag.Lookup(tag)
	if !ok {
		return f.Name, false
	}
	opts := strings.Split(val, ",")
	omit := false
	for _, key := range opts {
		if key == "omitempty" {
			omit = true
			break
		}
	}
	return opts[0], omit
}

// encode struct to map[string]inteface{}
func EncodeStructToMap(s interface{}, tag string) (res map[string]interface{}, err error) {
	res = map[string]interface{}{}
	v := reflect.ValueOf(s)
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not supported", t.Kind())
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// skip unexported fields. from godoc:
		// PkgPath is the package path that qualifies a lower case (unexported)
		// field name. It is empty for upper case (exported) field names.
		if f.PkgPath != "" {
			continue
		}

		fv := v.Field(i)
		key, omit := ReadTag(f, tag)
		// skip empty values when "omitempty" set.
		if omit && fv.String() == "" {
			continue
		}

		var val = fv.Interface()
		switch fv.Kind() {
		case reflect.Struct:
			val, err = EncodeStructToMap(fv.Interface(), tag)
			if err != nil {
				return nil, err
			}
		case reflect.Ptr:
			val, err = EncodeStructToMap(fv.Elem().Interface(), tag)
			if err != nil {
				return nil, err
			}
		default:

		}
		res[key] = val
	}
	return res, nil
}
