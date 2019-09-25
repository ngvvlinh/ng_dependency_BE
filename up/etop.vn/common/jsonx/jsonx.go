package jsonx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type modeType int

const (
	evaluating modeType = iota + 1
	safe
	revalidate
)

var mustValidate bool
var recognizedTypes = make(map[reflect.Type]modeType)
var m sync.RWMutex

func EnableValidation() {
	if mustValidate {
		panic("already enabled")
	}
	mustValidate = true
}

func Marshal(v interface{}) ([]byte, error) {
	if mustValidate {
		validate(v)
	}
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if mustValidate {
		validate(v)
	}
	return json.Unmarshal(data, v)
}

func validate(v interface{}) {
	value := reflect.Indirect(reflect.ValueOf(v))
	m.RLock()
	if recognizedTypes[value.Type()] == safe {
		m.RUnlock()
		return
	}
	m.RUnlock()
	m.Lock()
	defer m.Unlock()
	_, err := validateTag(value, nil)
	if err != nil {
		panic(err)
	}
}

func validateTag(v reflect.Value, t reflect.Type) (_mode modeType, _ error) {
	v = reflect.Indirect(v)
	// workaround for "call of reflect.Value.Type on zero Value"
	if t != nil {
		t = indirect(t)
	} else {
		t = v.Type()
	}

	// fast test
	currentMode := recognizedTypes[t]
	if currentMode == safe || currentMode == evaluating {
		return currentMode, nil
	}

	// temporary set to evaluating, set back to mode later
	_mode, recognizedTypes[t] = safe, evaluating
	defer func() {
		recognizedTypes[t] = _mode
	}()

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		elem := indirect(t.Elem())
		if recognizedTypes[elem] == safe {
			return safe, nil
		}
		if v.Len() == 0 {
			return validateTag(reflect.New(elem).Elem(), elem)
		}
		for i, n := 0, v.Len(); i < n; i++ {
			value := v.Index(i)
			mode, err := validateTag(value, elem)
			if err != nil {
				return 0, err
			}
			if mode > _mode {
				_mode = mode
			}
		}
		return _mode, nil

	case reflect.Map:
		elem := indirect(t.Elem())
		if recognizedTypes[elem] == safe {
			return safe, nil
		}
		if v.Len() == 0 {
			return validateTag(reflect.New(elem).Elem(), elem)
		}
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			mode, err := validateTag(value, elem)
			if err != nil {
				return 0, err
			}
			if mode > _mode {
				_mode = mode
			}
		}
		return _mode, nil

	case reflect.Interface:
		if v.IsNil() {
			return revalidate, nil
		}
		_, err := validateTag(v.Elem(), nil)
		return revalidate, err // always revalidate interface

	case reflect.Struct:
		// fast path: only validate field of type interface or revalidate
		if currentMode == revalidate {
			for i, n := 0, v.NumField(); i < n; i++ {
				tField := indirect(t.Field(i).Type)
				if recognizedTypes[tField] == safe {
					continue
				}
				_, err := validateTag(v.Field(i), nil)
				if err != nil {
					return 0, fmt.Errorf(
						"field %v of type %v: %v",
						t.Field(i).Name, t.Name(), err)
				}
			}
			return revalidate, nil
		}
		for i, n := 0, v.NumField(); i < n; i++ {
			vField := v.Field(i)
			tField := t.Field(i)
			jsonTag := tField.Tag.Get("json")
			if jsonTag == "" {
				return 0, fmt.Errorf(
					"field %v of type %v must have json tag",
					tField.Name, t.Name())
			}
			if jsonTag == "-" || strings.HasPrefix(jsonTag, "-,") {
				continue
			}
			mode, err := validateTag(vField, nil)
			if err != nil {
				return 0, fmt.Errorf(
					"field %v of type %v: %v",
					tField.Name, t.Name(), err)
			}
			if mode > _mode {
				_mode = mode
			}
		}
		return _mode, nil

	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return safe, nil

	case reflect.UnsafePointer,
		reflect.Func,
		reflect.Chan,
		reflect.Uintptr,
		reflect.Invalid:
		return 0, fmt.Errorf("invalid type %v for json", v)

	default:
		return 0, fmt.Errorf("unknown type %v for json", v)
	}
}

func indirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
