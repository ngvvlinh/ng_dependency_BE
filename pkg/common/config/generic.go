package cc

import (
	"reflect"

	"gopkg.in/yaml.v2"

	cm "o.o/backend/pkg/common"
	"o.o/common/l"
)

// GenericConfig reads config into a dynamic map. It's useful for registering and
// loading drivers at runtime. The empty value is ready for use.
//
// Currently, only yaml is supported for writing config.
type GenericConfig struct {
	m  map[string]*genericItem
	ms yaml.MapSlice

	// Mark the config as processed, which won't allow registering more keys.
	processed bool
}

type genericItem struct {
	v interface{} // store the config for decoding later

	// Mark the item as processed, which allows reporting missing keys.
	processed bool
}

// GenericConfig implements yaml.Unmarshaler
var _ yaml.Unmarshaler = (*GenericConfig)(nil)

// GenericConfig implements yaml.Unmarshaler
func (gc *GenericConfig) UnmarshalYAML(fn func(interface{}) error) error {
	gc.processed = true
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
		if gi != nil {
			gi.processed = true
			continue
		}

		// We ignore unknown keys, because we usually share config between
		// production and dev environments. When adding new drivers for dev,
		// it'll be not present for production, so we can safely ignore it.
		//
		// The reverse is also true for local development. We add a new key and
		// don't want to force everyone to add the key to their local config
		// file. They can just use the default value and focus on their job.
		ll.Warn("unknown config key (ignored)", l.String("key", key))
	}
	for k, v := range gc.m {
		if !v.processed {

			// Use the default value for missing keys. In the future, we may
			// allow configuring whether keys are required.
			ll.Warn("no config for key (used default)", l.String("key", k))
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
		if err = yaml.Unmarshal(out, gc.m[key].v); err != nil {
			return err
		}
	}
	return nil
}

// Register adds key and the variable for decoding config into. It may provide
// a default value.
//
// v must be a pointer.
func (gc *GenericConfig) Register(name string, v interface{}) {
	if gc.processed {
		ll.Panic("already processed")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		ll.S.Panicf("key %v: type %t must be a pointer", name, v)
	}
	if gc.m == nil {
		gc.m = make(map[string]*genericItem)
	}
	if gc.m[name] != nil {
		ll.Panic("already registered", l.String("name", name))
	}
	gc.m[name] = &genericItem{v: v}
}
