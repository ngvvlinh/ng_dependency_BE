package cc

import (
	"gopkg.in/yaml.v2"

	cm "o.o/backend/pkg/common"
	"o.o/common/l"
)

// GenericConfig read config into a dynamic map. It's useful for registering and
// loading drivers at runtime.
//
// Currently, only yaml is supported for writing config.
type GenericConfig struct {
	m  map[string]*GenericItem
	ms yaml.MapSlice
}

type GenericItem struct {
	v         interface{} // store the config for decoding later
	processed bool
}

// GenericConfig implements yaml.Unmarshaler
var _ yaml.Unmarshaler = (*GenericConfig)(nil)

// GenericConfig implements yaml.Unmarshaler
func (gc *GenericConfig) UnmarshalYAML(fn func(interface{}) error) error {
	if err := fn(&gc.ms); err != nil {
		return err
	}

	// first pass: validate registered keys and input
	for _, v := range gc.m {
		v.processed = false
	}
	for _, item := range gc.ms {
		key, ok := item.Key.(string)
		if !ok {
			return cm.Errorf(cm.Internal, nil, "unrecognized %+v (type %t)", key, key)
		}
		gi := gc.m[key]
		if gi == nil {
			return cm.Errorf(cm.Internal, nil, `unknown config key "%v"`, key)
		}
		gi.processed = true
	}
	for k, v := range gc.m {
		if !v.processed {
			return cm.Errorf(cm.Internal, nil, "no config for %v", k)
		}
	}

	// second pass: decode each item
	//
	// NOTE(vu): There may be a better way than Marshal then Unmarshal again,
	// but configs are usually small and loaded once at starting. So we chose
	// to keep the code simple.
	for _, item := range gc.ms {
		key := item.Key.(string)
		out, err := yaml.Marshal(item.Value)
		if err != nil {
			return cm.Errorf(cm.Internal, err, "can not marshal yaml")
		}
		if err = yaml.UnmarshalStrict(out, gc.m[key].v); err != nil {
			return err
		}
	}
	return nil
}

func (gc *GenericConfig) Register(name string, v interface{}) {
	if gc.m == nil {
		gc.m = make(map[string]*GenericItem)
	}
	if gc.m[name] != nil {
		ll.Panic("already register", l.String("name", name))
	}
	gc.m[name] = &GenericItem{v: v}
}
