package testing

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-test/deep"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
)

var ll = l.New()

// Spec ...
func Spec(items ...interface{}) {
	rename(items)
	convey.Convey(items...)
}

// SkipSpec ...
func SkipSpec(items ...interface{}) {
	rename(items)
	convey.SkipConvey(items...)
}

// FocusSpec ...
func FocusSpec(items ...interface{}) {
	rename(items)
	convey.FocusConvey(items...)
}

func rename(items []interface{}) {
	if len(items) > 0 {
		if s, ok := items[0].(string); ok {
			s += "\n"
			items[0] = s
		}
	}
}

// EqualSlice ...
func EqualSlice(t assert.TestingT, actual, expected interface{}, msgAndArgs ...interface{}) {
	const s0 = "Slice not equal: %v\n"
	const s1 = "  : %#v (expected)\n"
	const s2 = " != %#v (actual)"
	if msg := ShouldResembleSlice(actual, expected); msg != "" {
		assert.Fail(t, fmt.Sprintf(s1+s2, expected, actual), msgAndArgs...)
	}
}

// ShouldCMError ...
func ShouldCMError(actual interface{}, expected ...interface{}) string {
	return shouldError(actual, expected, func(actual, expected string) string {
		if actual != expected {
			return "message not equal"
		}
		return ""
	})
}

// ShouldCMErrorContains ...
func ShouldCMErrorContains(actual interface{}, expected ...interface{}) string {
	return shouldError(actual, expected, func(actual, expected string) string {
		if !strings.Contains(actual, expected) {
			return "not contain the message"
		}
		return ""
	})
}

func shouldError(actual interface{}, expected []interface{},
	fn func(actual, expected string) string) string {

	const msg = "Expected: %v '%v'\nActual:   %v '%v'\n(Should error: %v)!"

	if len(expected) != 2 {
		return fmt.Sprintf("This assertion requires exactly %v comparison values (you provided %v).", 2, len(expected))
	}

	code, ok := expected[0].(cm.Code)
	if !ok {
		return fmt.Sprintf("This assertion require the first comparison value is cm.Code (you provided %v)", reflect.TypeOf(expected[0]))
	}

	var s string
	if ss, ok := expected[1].(string); ok {
		s = ss
	} else if err, ok := expected[1].(error); ok {
		s = err.Error()
	} else {
		return fmt.Sprintf("This assertion requires the second comparison value as string or error (you provided %v).", expected)
	}
	if s == "" {
		s = cm.DefaultErrorMessage(code)
	}

	if actual == nil || reflect.ValueOf(actual).IsNil() {
		return fmt.Sprintf("expected error but got nil")
	}
	err, ok := actual.(error)
	if !ok {
		return fmt.Sprintf(msg,
			code, s,
			reflect.TypeOf(actual), actual,
			fmt.Sprintf("expected type error but got %v", reflect.TypeOf(actual)))
	}
	if cm.ErrorCode(err) != expected[0].(cm.Code) {
		return fmt.Sprintf(msg,
			code, s,
			cm.ErrorCode(err), actual,
			"message not equal")
	}
	if errMsg := fn(err.Error(), s); errMsg != "" {
		return fmt.Sprintf(msg,
			code, s,
			cm.ErrorCode(err), actual,
			errMsg)
	}
	return ""
}

// ShouldResembleSlice ...
func ShouldResembleSlice(actual interface{}, expected ...interface{}) string {
	const msg = "Expected: '%v'\nActual:   '%v'\n(Should equal slice: %v)!"

	if len(expected) != 1 {
		return fmt.Sprintf("This assertion requires exactly %v comparison values (you provided %v).", 1, len(expected))
	}

	a := reflect.ValueOf(actual)
	e := reflect.ValueOf(expected[0])
	if a.Kind() != reflect.Slice || e.Kind() != reflect.Slice {
		return fmt.Sprintf(msg, l.Dump(expected[0]), l.Dump(actual),
			"Both must be slice")
	}

	if a.Len() != e.Len() {
		return fmt.Sprintf(msg, l.Dump(expected[0]), l.Dump(actual),
			"Length not equal")
	}

	count := 0
	indexes := make([]bool, a.Len())
	for i := 0; i < a.Len(); i++ {
		ai := a.Index(i)
		matched := false
		for j := 0; j < e.Len(); j++ {
			if ok := indexes[j]; ok {
				// already compared
				continue
			}
			ej := e.Index(j)
			if res := deep.Equal(ej.Interface(), ai.Interface()); res == nil {
				indexes[j] = true
				count++
				matched = true
				break
			}
		}
		if !matched {
			diff := deep.Equal(expected[0], actual)
			format := "Not match %d items: %v"
			if len(diff) == 1 {
				format = "Not match %d item: %v"
			}
			return fmt.Sprintf(msg, l.Dump(expected[0]), l.Dump(actual),
				fmt.Sprintf(format, len(diff), l.Dump(diff)),
			)
		}
	}
	if count != a.Len() {
		return fmt.Sprintf(msg, l.Dump(expected[0]), l.Dump(actual),
			"Slices are not equal")
	}
	return ""
}

// ShouldDeepEqual ...
func ShouldDeepEqual(actual interface{}, expected ...interface{}) string {
	diff := deep.Equal(expected[0], actual)
	if len(diff) == 0 {
		return ""
	}

	format := "Not match %d items: %v"
	if len(diff) == 1 {
		format = "Not match %d item: %v"
	}

	const msg = "Expected: '%v'\nActual:   '%v'\n(Should deep equal: %v)!"
	return fmt.Sprintf(msg, l.Dump(expected[0]), l.Dump(actual),
		fmt.Sprintf(format, len(diff), l.Dump(diff)),
	)
}

// ShouldResembleByKey ...
func ShouldResembleByKey(key string) func(actual interface{}, expected ...interface{}) string {
	const msg = "Expected: '%v'\nActual:   '%v'\n(Should equal slice: %v)!"

	return func(actual interface{}, expected ...interface{}) string {
		if len(expected) != 1 {
			return fmt.Sprintf("This assertion requires exactly %v comparison values (you provided %v).", 1, len(expected))
		}
		formatError := func(format string, args ...interface{}) string {
			errMsg := fmt.Sprintf(format, args...)
			return fmt.Sprintf(msg, l.Dump(expected[0]), l.Dump(actual), errMsg)
		}

		a := reflect.ValueOf(actual)
		e := reflect.ValueOf(expected[0])
		if a.Kind() != reflect.Slice || e.Kind() != reflect.Slice {
			return formatError("Both must be slice")
		}
		if errMsg := canGetKey(a.Type().Elem(), key); errMsg != "" {
			return formatError(errMsg)
		}
		if errMsg := canGetKey(e.Type().Elem(), key); errMsg != "" {
			return formatError(errMsg)
		}
		if a.Len() != e.Len() {
			return formatError("Length not equal")
		}

		collectIndexes := func(name string, list reflect.Value) (
			[]reflect.Value, map[interface{}]int, string,
		) {
			keys := make([]reflect.Value, list.Len())
			mapIndexes := make(map[interface{}]int)
			for i := 0; i < list.Len(); i++ {
				item := list.Index(i)
				switch item.Kind() {
				case reflect.Interface, reflect.Map, reflect.Ptr:
					if item.IsNil() {
						return nil, nil, formatError(
							"All items must not be nil (%v[%v] is nil)", name, i)
					}
				}

				itemKey := getKey(item, key)
				keys[i] = itemKey
				if !itemKey.IsValid() {
					return nil, nil, formatError(
						"Could not get key from %v[%v]", name, i)
				}

				keyValue := itemKey.Interface()
				if keyValue == nil {
					return nil, nil, formatError(
						"All item keys must not be nil (%v[%v].%v is nil)", name, i, key)
				}

				keyType := reflect.TypeOf(keyValue)
				if !keyType.Comparable() {
					return nil, nil, formatError(
						"All item keys must be comparable (%v[%v].%v is not, type is `%v`)",
						name, i, key, keyType)
				}
				if prev, ok := mapIndexes[keyValue]; ok {
					return nil, nil, formatError(
						"%v[%v] and %v[%v] has duplicated keys: `%v`",
						name, prev, name, i, keyValue)
				}
				mapIndexes[keyValue] = i
			}
			return keys, mapIndexes, ""
		}
		expectedKeys, _, errMsg := collectIndexes("expected", e)
		if errMsg != "" {
			return errMsg
		}
		_, mapActualIndexes, errMsg := collectIndexes("actual", a)
		if errMsg != "" {
			return errMsg
		}

		// Compare actual with the same order as expected
		for i, ekey := range expectedKeys {
			actualIndex, ok := mapActualIndexes[ekey.Interface()]
			if !ok {
				return formatError("Expected item with %v=`%v` but not found",
					key, ekey.Interface())
			}

			expectedItem := e.Index(i)
			actualItem := a.Index(actualIndex)

			diff := deep.Equal(actualItem.Interface(), expectedItem.Interface())
			if len(diff) > 0 {
				return formatError("Item with %v=`%v` is different: %v",
					key, ekey.Interface(), l.Dump(diff))
			}
		}
		return ""
	}
}

func canGetKey(t reflect.Type, key string) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Map, reflect.Interface:
		return ""
	case reflect.Struct:
		_, ok := t.FieldByName(key)
		if !ok {
			similar := ""
			for i, n := 0, t.NumField(); i < n; i++ {
				name := t.Field(i).Name
				if strings.ToLower(name) == strings.ToLower(key) {
					similar = name
				}
			}
			if similar != "" {
				return fmt.Sprintf(
					"Key `%v` not found in struct (but it has `%v`)",
					key, similar)
			}
			return fmt.Sprintf("Key `%v` not found in struct", key)
		}
		return ""
	}
	return "Both must be slice of struct, *struct, map or interface"
}

func getKey(v reflect.Value, key string) reflect.Value {
	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
		if !v.IsValid() {
			return v
		}
	}
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return v
	}

	switch v.Kind() {
	case reflect.Map:
		return v.MapIndex(reflect.ValueOf(key))
	case reflect.Struct:
		return v.FieldByName(key)
	default:
		return reflect.Value{}
	}
}
