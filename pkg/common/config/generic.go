package cc

import (
	"context"
	"reflect"
	"sort"

	"gopkg.in/yaml.v2"

	cm "o.o/backend/pkg/common"
	"o.o/common/l"
)

// GenericConfig reads config into a dynamic map. It's useful for registering
// and loading drivers at runtime. The empty value is ready for use.
//
// Currently, only yaml is supported for writing config.
type GenericConfig struct {
	m map[string]*genericItem

	// Mark the config as init, which won't allow registering more keys.
	init bool
}

type genericItem struct {
	v interface{} // store the config for decoding later

	// Mark the item as processed, which allows reporting missing keys.
	processed bool
}

// GenericConfig implements yaml.Unmarshaler
var _ yaml.Unmarshaler = (*GenericConfig)(nil)

func (gc *GenericConfig) MarshalYAML() (interface{}, error) {
	gc.init = true

	items := make(yaml.MapSlice, 0, len(gc.m))
	for k, v := range gc.m {
		items = append(items, yaml.MapItem{Key: k, Value: v})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Key.(string) < items[j].Key.(string) })
	return yaml.Marshal(items)
}

// GenericConfig implements yaml.Unmarshaler
func (gc *GenericConfig) UnmarshalYAML(fn func(interface{}) error) error {
	gc.init = true

	var items yaml.MapSlice
	if err := fn(&items); err != nil {
		return err
	}

	// first pass: validate registered keys and input
	for _, v := range gc.m {
		v.processed = false
	}
	for _, item := range items {
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
	for _, item := range items {
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
func (gc *GenericConfig) Register(name string, cfg, constructor interface{}) {
	if gc.init {
		ll.Panic("already processed")
	}
	if reflect.TypeOf(cfg).Kind() != reflect.Ptr {
		ll.S.Panicf("key %v: type %t must be a pointer", name, cfg)
	}
	if gc.m == nil {
		gc.m = make(map[string]*genericItem)
	}
	if gc.m[name] != nil {
		ll.Panic("already registered", l.String("name", name))
	}
	gc.m[name] = &genericItem{v: cfg}
}

func (gc *GenericConfig) Build(ctx context.Context, result interface{}) error {
	return nil
}
