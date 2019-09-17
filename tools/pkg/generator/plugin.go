package generator

import "go/types"

type Filter interface {
	FilterPackage(*PreparsedPackage) (bool, error)
}

type Qualifier interface {
	Qualify(*types.Package) string
}

type Plugin interface {
	Name() string
	FilterPackage(*PreparsedPackage) (bool, error)
	Generate(Engine) error
}

type pluginStruct struct {
	name      string
	plugin    Plugin
	enabled   bool
	qualifier types.Qualifier

	includes  []bool
	includesN int
}

func RegisterPlugin(plugins ...Plugin) error {
	for _, plugin := range plugins {
		if err := theEngine.registerPlugin(plugin); err != nil {
			return Errorf(err, "register plugin %v: %v", plugin.Name(), err)
		}
	}
	return nil
}

func (ng *engine) registerPlugin(plugin Plugin) error {
	name := plugin.Name()
	if name == "" {
		return Errorf(nil, "empty name")
	}
	if plugin == nil {
		return Errorf(nil, "nil plugin")
	}
	if ng.pluginsMap[name] != nil {
		return Errorf(nil, "duplicated pluginStruct name: %v", name)
	}

	pl := &pluginStruct{name: name, plugin: plugin}
	if q, ok := plugin.(Qualifier); ok {
		pl.qualifier = q.Qualify
	}

	ng.plugins = append(ng.plugins, pl)
	ng.pluginsMap[name] = pl
	return nil
}
