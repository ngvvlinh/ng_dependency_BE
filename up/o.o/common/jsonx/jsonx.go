package jsonx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"sync"
)

const DefaultRoute = "/==/jsonx"

func RegisterHTTPHandler(mux *http.ServeMux) {
	mux.HandleFunc(DefaultRoute, ServeHTTP)
}

type ErrorMode int

const (
	Warning ErrorMode = iota + 1
	Panicking
)

type modeType int

const (
	evaluating modeType = iota + 1
	safe
	revalidate
)

type routeType int

const (
	marshal routeType = iota + 1
	unmarshal
)

var marshaler = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
var unmarshaler = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()

var enabledMode ErrorMode
var recognizedTypes = make(map[reflect.Type]modeType)
var savedErrors = make(map[Error]struct{})
var mr, me sync.RWMutex

type Error struct {
	Message string
}

func GetErrors() []Error {
	me.RLock()
	defer me.RUnlock()

	if len(savedErrors) == 0 {
		return nil
	}
	errs := make([]Error, 0, len(savedErrors))
	for msg := range savedErrors {
		errs = append(errs, msg)
	}
	sort.Slice(errs, func(i, j int) bool {
		return errs[i].Message < errs[j].Message
	})
	return errs
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	errs := GetErrors()
	if len(errs) == 0 {
		_, _ = fmt.Fprint(w, "no error")
		return
	}
	for _, err := range errs {
		_, _ = fmt.Fprint(w, err.Message, "\n")
	}
}

func EnableValidation(mode ErrorMode) {
	if mode <= 0 {
		panic("invalid mode")
	}
	if enabledMode > 0 {
		panic("already enabled")
	}
	enabledMode = mode
}

func Marshal(v interface{}) ([]byte, error) {
	if enabledMode > 0 {
		mustValidate(v, marshal)
	}
	workaroundFillSlice(v)
	return json.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	var b strings.Builder
	err := MarshalTo(&b, v)
	if err != nil {
		return "", nil
	}
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func MustMarshalToString(v interface{}) string {
	s, err := MarshalToString(v)
	if err != nil {
		panic(err)
	}
	return s
}

func MarshalTo(w io.Writer, v interface{}) error {
	if enabledMode > 0 {
		mustValidate(v, marshal)
	}
	workaroundFillSlice(v)
	return json.NewEncoder(w).Encode(v)
}

func Unmarshal(data []byte, v interface{}) error {
	if enabledMode > 0 {
		mustValidate(v, unmarshal)
	}
	return json.Unmarshal(data, v)
}

func UnmarshalString(s string, v interface{}) error {
	return Unmarshal([]byte(s), v)
}

func UnmarshalFrom(r io.Reader, v interface{}) error {
	if enabledMode > 0 {
		mustValidate(v, unmarshal)
	}
	return json.NewDecoder(r).Decode(v)
}

func mustValidate(v interface{}, route routeType) {
	_, _, err := validate(v, route)
	if err == nil {
		return
	}
	if enabledMode == Panicking {
		panic(err)
	}

	// save the error
	msg := Error{Message: err.Error()}
	me.RLock()
	if _, ok := savedErrors[msg]; ok {
		me.RUnlock()
		return
	}
	me.RUnlock()
	me.Lock()
	defer me.Unlock()

	savedErrors[msg] = struct{}{}
}

func validate(v interface{}, route routeType) (fastpath int, _ modeType, _ error) {
	t := indirectType(reflect.TypeOf(v))
	if t == nil {
		return 1, safe, nil
	}
	mr.RLock()
	if recognizedTypes[t] == safe {
		fastpath = 2
		mr.RUnlock()
		return 2, safe, nil
	}
	mr.RUnlock()
	mr.Lock()
	defer mr.Unlock()
	mode, err := validateTag(reflect.ValueOf(v), t, route)
	return 0, mode, err
}

func validateTag(v reflect.Value, t reflect.Type, route routeType) (_mode modeType, _ error) {
	v = indirectValue(v)
	t = indirectType(t)
	if t == nil {
		return safe, nil
	}

	// fast test
	currentMode := recognizedTypes[t]
	if currentMode == safe || currentMode == evaluating {
		return currentMode, nil
	}
	if v.Kind() == reflect.Invalid {
		return validateTag(reflect.New(t).Elem(), t, route)
	}

	// temporary set to evaluating, set back to mode later
	defer func() { recognizedTypes[t] = _mode }()
	_mode, recognizedTypes[t] = safe, evaluating

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		elem := indirectType(t.Elem())
		if recognizedTypes[elem] == safe {
			return safe, nil
		}
		if v.Len() == 0 {
			return validateTag(reflect.New(elem).Elem(), elem, route)
		}
		for i, n := 0, v.Len(); i < n; i++ {
			value := v.Index(i)
			mode, err := validateTag(value, elem, route)
			if err != nil {
				return 0, err
			}
			if mode > _mode {
				_mode = mode
			}
		}
		return _mode, nil

	case reflect.Map:
		elem := indirectType(t.Elem())
		if recognizedTypes[elem] == safe {
			return safe, nil
		}
		if v.Len() == 0 {
			return validateTag(reflect.New(elem).Elem(), elem, route)
		}
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			mode, err := validateTag(value, elem, route)
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
		_, err := validateTag(v.Elem(), v.Elem().Type(), route)
		return revalidate, err // always revalidate interface

	case reflect.Struct:
		// fast path: only validate field of type interface or revalidate
		if currentMode == revalidate {
			for i, n := 0, v.NumField(); i < n; i++ {
				tField := indirectType(t.Field(i).Type)
				if recognizedTypes[tField] == safe {
					continue
				}
				_, err := validateTag(v.Field(i), tField, route)
				if err != nil {
					return 0, fmt.Errorf(
						"field %v of type %v.%v: %v",
						t.Field(i).Name, t.PkgPath(), t.Name(), err)
				}
			}
			return revalidate, nil
		}
		for i, n := 0, v.NumField(); i < n; i++ {
			vField := v.Field(i)
			tField := t.Field(i)
			jsonTag := tField.Tag.Get("json")
			if jsonTag == "" && !tField.Anonymous {
				return 0, fmt.Errorf(
					"field %v of type %v.%v must have json tag",
					tField.Name, t.PkgPath(), t.Name())
			}
			if jsonTag == "-" || strings.HasPrefix(jsonTag, "-,") {
				continue
			}
			switch route {
			case marshal:
				if vField.Type().Implements(marshaler) {
					continue
				}
			case unmarshal:
				if vField.Type().Implements(unmarshaler) || reflect.New(vField.Type()).Type().Implements(unmarshaler) {
					continue
				}
			default:
				panic("unexpected")
			}
			mode, err := validateTag(vField, tField.Type, route)
			if err != nil {
				return 0, fmt.Errorf(
					"field %v of type %v.%v: %v",
					tField.Name, t.PkgPath(), t.Name(), err)
			}
			if mode > _mode {
				_mode = mode
			}
		}
		return _mode, nil

	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
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

func indirectValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func indirectType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func workaroundFillSlice(v interface{}) {
	val := reflect.ValueOf(v)
	workaroundFillSlice0(val)
}

func workaroundFillSlice0(val reflect.Value) {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	typ := val.Type()
	for i, n := 0, val.NumField(); i < n; i++ {
		field := reflect.Indirect(val.Field(i))
		fieldType := typ.Field(i)
		if fieldType.Tag == "" {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			workaroundFillSlice0(field)

		case reflect.Slice:
			if field.IsNil() {
				// workaround: set it to empty slice ([]Type)
				// TODO: remove workaround
				emptySlice := reflect.MakeSlice(field.Type(), 0, 0)
				field.Set(emptySlice)
				continue
			}
			elemType := fieldType.Type.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			if elemType.Kind() == reflect.Struct && elemType.Name() != "Time" {
				for i, n := 0, field.Len(); i < n; i++ {
					workaroundFillSlice0(field.Index(i))
				}
			}
		}
	}
}
