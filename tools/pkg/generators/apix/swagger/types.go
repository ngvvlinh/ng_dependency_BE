package swagger

import (
	"fmt"
	"go/types"

	"etop.vn/backend/tools/pkg/generator"
)

var (
	typeID        types.Type
	typeTime      types.Type
	typeTimestamp types.Type
)

func initTypes(ng generator.Engine) {
	populateType(ng, &typeID, "etop.vn/capi/dot", "ID")
	populateType(ng, &typeTime, "time", "Time")
	populateType(ng, &typeTimestamp, "etop.vn/capi/dot", "Time")
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

func isTime(typ types.Type) bool {
	typ = skipPointer(typ)
	return typ == typeTime || typ == typeTimestamp
}

func skipPointer(typ types.Type) types.Type {
	ptr, ok := typ.(*types.Pointer)
	if ok {
		return ptr.Elem()
	}
	return typ
}

func isNamedStruct(typ types.Type, inner *types.Type) bool {
	typ = skipPointer(typ)
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	st, ok := named.Underlying().(*types.Struct)
	*inner = st
	return ok
}

func isArray(typ types.Type, inner *types.Type) bool {
	slice, ok := typ.(*types.Slice)
	if ok {
		*inner = slice.Elem()
	}
	return ok
}

func isSliceOfBytes(typ types.Type) bool {
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

func isMap(typ types.Type) bool {
	_, ok := typ.(*types.Map)
	return ok
}

func isEnum(typ types.Type, inner *types.Type) bool {
	typ = skipPointer(typ)
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	basic, ok := named.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	ok = false
	for i, n := 0, named.NumMethods(); i < n; i++ {
		method := named.Method(i)
		if method.Name() == "Enum" {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	*inner = basic
	return true
}

func isID(typ types.Type) bool {
	typ = skipPointer(typ)
	return typ == typeID
}

func isBasic(typ types.Type, inner *types.Type) bool {
	typ = skipPointer(typ)
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	*inner = basic
	return true
}

func isNamedInterface(typ types.Type, inner *types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	iface, ok := named.Underlying().(*types.Interface)
	*inner = iface
	return ok
}
