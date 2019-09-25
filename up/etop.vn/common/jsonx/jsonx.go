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

var fastpath int // for testing only
var enabled bool
var recognizedTypes = make(map[reflect.Type]modeType)
var m sync.RWMutex

func EnableValidation() {
	if enabled {
		panic("already enabled")
	}
	enabled = true
}

func Marshal(v interface{}) ([]byte, error) {
	if enabled {
		mustValidate(v)
	}
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if enabled {
		mustValidate(v)
	}
	return json.Unmarshal(data, v)
}

func mustValidate(v interface{}) {
	_, err := validate(v)
	if err != nil {
		panic(err)
	}
}

func validate(v interface{}) (modeType, error) {
	fastpath = 0
	t := indirect(reflect.TypeOf(v))
	if t == nil {
		fastpath = 1
		return safe, nil
	}
	m.RLock()
	if recognizedTypes[t] == safe {
		fastpath = 2
		m.RUnlock()
		return safe, nil
	}
	m.RUnlock()
	m.Lock()
	defer m.Unlock()
	return validateTag(reflect.ValueOf(v), t)
}

func validateTag(v reflect.Value, t reflect.Type) (_mode modeType, _ error) {
	defer func() { recognizedTypes[t] = _mode }()

	v = reflect.Indirect(v)
	t = indirect(t)
	if t == nil {
		return safe, nil
	}

	// fast test
	currentMode := recognizedTypes[t]
	if currentMode == safe || currentMode == evaluating {
		return currentMode, nil
	}
	if v.Kind() == reflect.Invalid {
		return validateTag(reflect.New(t).Elem(), t)
	}

	// temporary set to evaluating, set back to mode later
	_mode, recognizedTypes[t] = safe, evaluating

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
		_, err := validateTag(v.Elem(), v.Elem().Type())
		return revalidate, err // always revalidate interface

	case reflect.Struct:
		// fast path: only validate field of type interface or revalidate
		if currentMode == revalidate {
			for i, n := 0, v.NumField(); i < n; i++ {
				tField := indirect(t.Field(i).Type)
				if recognizedTypes[tField] == safe {
					continue
				}
				_, err := validateTag(v.Field(i), tField)
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
			mode, err := validateTag(vField, tField.Type)
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
	if t == nil {
		return nil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
