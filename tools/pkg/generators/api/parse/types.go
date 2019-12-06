package parse

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/api/defs"
)

type MetaKey struct {
	Item interface{}
	Key  interface{}
}

type Meta map[MetaKey]interface{}

func (m Meta) Get(item, key interface{}) interface{} {
	return m[MetaKey{item, key}]
}

func (m Meta) Set(item, key, value interface{}) {
	m[MetaKey{item, key}] = value
}

func (m Meta) Cache(item, key interface{}, fn func() interface{}) interface{} {
	if value := m.Get(item, key); value != nil {
		return value
	}
	value := fn()
	m.Set(item, key, value)
	return value
}

const dotPkgPath = "etop.vn/capi/dot"

type Info struct {
	Meta
	ng generator.Engine

	typeError   types.Type
	typeStdTime types.Type
	typeDotTime types.Type
	typeID      types.Type
	typeNullID  types.Type
}

func NewInfo(ng generator.Engine) *Info {
	inf := &Info{Meta: make(Meta), ng: ng}
	populateType(ng, &inf.typeError, "", "error")
	populateType(ng, &inf.typeStdTime, "time", "Time")
	populateType(ng, &inf.typeDotTime, dotPkgPath, "Time")
	populateType(ng, &inf.typeID, dotPkgPath, "ID")
	populateType(ng, &inf.typeNullID, dotPkgPath, "NullID")
	return inf
}

func populateType(ng generator.Engine, typ *types.Type, pkgPath, name string) {
	if pkgPath == "" {
		*typ = ng.GetBuiltinType(name)
		return
	}

	obj := ng.GetObjectByName(pkgPath, name)
	if obj == nil {
		panic(fmt.Sprintf("type %v.%v not found", pkgPath, name))
	}
	*typ = obj.Type()
}

func (inf *Info) IsTime(typ types.Type) bool {
	typ = SkipPointer(typ)
	return typ == inf.typeStdTime || typ == inf.typeDotTime
}

func SkipPointer(typ types.Type) types.Type {
	ptr, ok := typ.(*types.Pointer)
	if ok {
		return ptr.Elem()
	}
	return typ
}

func (inf *Info) IsNullID(typ types.Type) bool {
	return typ == inf.typeNullID
}

func (inf *Info) IsNullBasic(typ types.Type, inner *types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	pkg := named.Obj().Pkg()
	if pkg == nil {
		return false
	}
	if pkg.Path() != dotPkgPath {
		return false
	}
	name := named.Obj().Name()
	if strings.HasPrefix(name, "Null") {
		field := named.Underlying().(*types.Struct).Field(0)
		if field.Name() == "Valid" {
			panic(fmt.Sprintf("invalid type %v", named))
		}
		*inner = field.Type()
		return true
	}
	return false
}

func (inf *Info) IsBasic(typ types.Type, inner *types.Type) bool {
	typ = SkipPointer(typ)
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	*inner = basic
	return true
}

func (inf *Info) IsNamedStruct(typ types.Type, inner *types.Type) bool {
	typ = SkipPointer(typ)
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	st, ok := named.Underlying().(*types.Struct)
	*inner = st
	return ok
}

func (inf *Info) IsArray(typ types.Type, inner *types.Type) bool {
	slice, ok := typ.(*types.Slice)
	if ok {
		*inner = slice.Elem()
	}
	return ok
}

func (inf *Info) IsSliceOfBytes(typ types.Type) bool {
	slice, ok := typ.(*types.Slice)
	if !ok {
		return false
	}
	basic, ok := slice.Elem().(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.Byte
}

func (inf *Info) IsMap(typ types.Type) bool {
	_, ok := typ.(*types.Map)
	return ok
}

func (inf *Info) IsEnum(typ types.Type) bool {
	typ = SkipPointer(typ)
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	_, ok = named.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return inf.GetEnum(typ) != nil
}

type keyEnum struct{}

func (inf *Info) GetEnum(typ types.Type) *defs.Enum {
	typ = SkipPointer(typ)
	return inf.Cache(typ, keyEnum{}, func() interface{} {
		obj := typ.(*types.Named).Obj()
		pkgPath := obj.Pkg().Path()
		pkg := inf.ng.GetPackageByPath(pkgPath)
		if pkg == nil {
			panic(fmt.Sprintf("package %v not found", pkgPath))
		}
		mapEnum := inf.parseEnumInPackage(pkg)
		return mapEnum[obj.Name()]
	}).(*defs.Enum)
}

func (inf *Info) parseEnumInPackage(pkg *packages.Package) map[string]*defs.Enum {
	return inf.Cache(pkg, keyEnum{}, func() interface{} {
		mapEnum, err := parseEnumInPackage(inf.ng, pkg)
		if err != nil {
			fmt.Printf("%+v", err)
			panic(fmt.Sprintf("can not parse enum in package %v", pkg.PkgPath))
		}
		return mapEnum
	}).(map[string]*defs.Enum)
}

func (inf *Info) IsID(typ types.Type) bool {
	typ = SkipPointer(typ)
	return typ == inf.typeID
}

func (inf *Info) IsNamedInterface(typ types.Type, inner *types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	iface, ok := named.Underlying().(*types.Interface)
	*inner = iface
	return ok
}
