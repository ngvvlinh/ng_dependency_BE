package scheme

import (
	"errors"
	"fmt"
	"reflect"
)

var Global = NewScheme()

func Convert(arg, out interface{}) error {
	return Global.Convert(arg, out)
}

type TypePair struct {
	Slice bool
	Arg   reflect.Type
	Out   reflect.Type
}

type ConversionFunc func(arg, out interface{}) error

type Scheme struct {
	convPairs map[TypePair]ConversionFunc
	ready     bool
}

func NewScheme() *Scheme {
	return &Scheme{
		convPairs: make(map[TypePair]ConversionFunc),
	}
}

func (s *Scheme) Register(arg, out interface{}, fn ConversionFunc) {
	if s.ready {
		panic("register too late!")
	}
	pair, err := getTypePair(arg, out)
	if err != nil {
		panic(err)
	}
	s.convPairs[pair] = fn
}

func (s *Scheme) Convert(arg, out interface{}) error {
	if !s.ready {
		s.ready = true
	}
	pair, err := getTypePair(arg, out)
	if err != nil {
		panic(fmt.Sprintf("invalid conversion type pair: %v (%T and %T)", err, arg, out))
	}
	fn := s.convPairs[pair]
	if fn == nil {
		panic(fmt.Sprintf("no conversion between %t and %t", arg, out))
	}
	return fn(arg, out)
}

func getTypePair(arg, out interface{}) (TypePair, error) {
	argType := reflect.TypeOf(arg)
	outType := reflect.TypeOf(out)
	switch {
	case argType.Kind() == reflect.Slice && outType.Kind() == reflect.Slice:
		return TypePair{}, errors.New("second param must be pointer to slice")

	case argType.Kind() == reflect.Slice &&
		outType.Kind() == reflect.Ptr && outType.Elem().Kind() == reflect.Slice:
		if argType.Elem().Kind() == reflect.Ptr &&
			outType.Elem().Elem().Kind() == reflect.Ptr {
			pair := TypePair{
				Slice: true,
				Arg:   argType.Elem().Elem(),
				Out:   outType.Elem().Elem().Elem(),
			}
			return pair, nil
		}
		return TypePair{}, errors.New("must be slice of pointer")

	case argType.Kind() != reflect.Slice && outType.Kind() != reflect.Slice:
		if argType.Kind() == reflect.Ptr &&
			outType.Kind() == reflect.Ptr {
			pair := TypePair{
				Slice: false,
				Arg:   argType.Elem(),
				Out:   outType.Elem(),
			}
			return pair, nil
		}
		return TypePair{}, errors.New("must be pointer")

	default:
		return TypePair{}, errors.New("both types must match")
	}
}
