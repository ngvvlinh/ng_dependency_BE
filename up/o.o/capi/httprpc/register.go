package httprpc

import (
	"reflect"

	"o.o/common/xerrors"
)

var globalRegistry = Registry{}

type RegisterFunc func(builder interface{}, hooks ...*Hooks) (Server, bool)

type Registry struct {
	funcs []RegisterFunc
}

func (r *Registry) Register(fn RegisterFunc) {
	r.funcs = append(r.funcs, fn)
}

func (r *Registry) NewServer(builder interface{}, hooks ...*Hooks) (Server, error) {
	for _, fn := range r.funcs {
		if server, _ := fn(builder, hooks...); server != nil {
			return server, nil
		}
	}
	if reflect.TypeOf(builder).Kind() != reflect.Func {
		return nil, xerrors.Errorf(xerrors.Internal, nil, "builder of type %T is not a function", builder)
	}
	return nil, xerrors.Errorf(xerrors.Internal, nil, "builder of type %T is not recognized", builder)
}

func (r *Registry) NewServers(builders []interface{}, hooks ...*Hooks) (servers []Server, _ error) {
	for _, builder := range builders {
		server, err := r.NewServer(builder, hooks...)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func Register(fn RegisterFunc) {
	globalRegistry.Register(fn)
}

func NewServer(builder interface{}, hooks ...*Hooks) (Server, error) {
	return globalRegistry.NewServer(builder, hooks...)
}

func MustNewServer(builder interface{}, hooks ...*Hooks) Server {
	server, err := globalRegistry.NewServer(builder, hooks...)
	must(err)
	return server
}

func NewServers(builders []interface{}, hooks ...*Hooks) (servers []Server, _ error) {
	return globalRegistry.NewServers(builders, hooks...)
}

func MustNewServers(builders []interface{}, hooks ...*Hooks) (servers []Server) {
	servers, err := globalRegistry.NewServers(builders, hooks...)
	must(err)
	return servers
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
