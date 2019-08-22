package generator

type Filter interface {
	FilterPackage(*PreparsedPackage) (bool, error)
}

type Plugin interface {
	Name() string
	FilterPackage(*PreparsedPackage) (bool, error)
	Generate(Engine) error
}

type pluginStruct struct {
	name    string
	plugin  Plugin
	enabled bool

	includes  []bool
	includesN int
}

func RegisterPlugin(plugins ...Plugin) error {
	for _, plugin := range plugins {
		if err := theEngine.registerPlugin(plugin); err != nil {
			return errorf(err, "register plugin %v: %v", plugin.Name(), err)
		}
	}
	return nil
}

func (ng *engine) registerPlugin(plugin Plugin) error {
	name := plugin.Name()
	if name == "" {
		return errorf(nil, "empty name")
	}
	if plugin == nil {
		return errorf(nil, "nil plugin")
	}
	if ng.pluginsMap[name] != nil {
		return errorf(nil, "duplicated pluginStruct name: %v", name)
	}

	pl := &pluginStruct{name: name, plugin: plugin}
	ng.plugins = append(ng.plugins, pl)
	ng.pluginsMap[name] = pl
	return nil
}
