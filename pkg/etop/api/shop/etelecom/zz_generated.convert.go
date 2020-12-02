// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package etelecom

import (
	time "time"

	etelecom "o.o/api/etelecom"
	shop "o.o/api/top/int/shop"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*etelecom.CreateExtensionArgs)(nil), (*shop.Extension)(nil), func(arg, out interface{}) error {
		Apply_etelecom_CreateExtensionArgs_shop_Extension(arg.(*etelecom.CreateExtensionArgs), out.(*shop.Extension))
		return nil
	})
	s.Register((*etelecom.Extension)(nil), (*shop.Extension)(nil), func(arg, out interface{}) error {
		Convert_etelecom_Extension_shop_Extension(arg.(*etelecom.Extension), out.(*shop.Extension))
		return nil
	})
	s.Register(([]*etelecom.Extension)(nil), (*[]*shop.Extension)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_Extensions_shop_Extensions(arg.([]*etelecom.Extension))
		*out.(*[]*shop.Extension) = out0
		return nil
	})
	s.Register((*shop.Extension)(nil), (*etelecom.Extension)(nil), func(arg, out interface{}) error {
		Convert_shop_Extension_etelecom_Extension(arg.(*shop.Extension), out.(*etelecom.Extension))
		return nil
	})
	s.Register(([]*shop.Extension)(nil), (*[]*etelecom.Extension)(nil), func(arg, out interface{}) error {
		out0 := Convert_shop_Extensions_etelecom_Extensions(arg.([]*shop.Extension))
		*out.(*[]*etelecom.Extension) = out0
		return nil
	})
	s.Register((*etelecom.ExtensionExternalData)(nil), (*shop.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		Convert_etelecom_ExtensionExternalData_shop_ExtensionExternalData(arg.(*etelecom.ExtensionExternalData), out.(*shop.ExtensionExternalData))
		return nil
	})
	s.Register(([]*etelecom.ExtensionExternalData)(nil), (*[]*shop.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_etelecom_ExtensionExternalDatas_shop_ExtensionExternalDatas(arg.([]*etelecom.ExtensionExternalData))
		*out.(*[]*shop.ExtensionExternalData) = out0
		return nil
	})
	s.Register((*shop.ExtensionExternalData)(nil), (*etelecom.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		Convert_shop_ExtensionExternalData_etelecom_ExtensionExternalData(arg.(*shop.ExtensionExternalData), out.(*etelecom.ExtensionExternalData))
		return nil
	})
	s.Register(([]*shop.ExtensionExternalData)(nil), (*[]*etelecom.ExtensionExternalData)(nil), func(arg, out interface{}) error {
		out0 := Convert_shop_ExtensionExternalDatas_etelecom_ExtensionExternalDatas(arg.([]*shop.ExtensionExternalData))
		*out.(*[]*etelecom.ExtensionExternalData) = out0
		return nil
	})
}

//-- convert o.o/api/top/int/shop.Extension --//

func Apply_etelecom_CreateExtensionArgs_shop_Extension(arg *etelecom.CreateExtensionArgs, out *shop.Extension) *shop.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shop.Extension{}
	}
	apply_etelecom_CreateExtensionArgs_shop_Extension(arg, out)
	return out
}

func apply_etelecom_CreateExtensionArgs_shop_Extension(arg *etelecom.CreateExtensionArgs, out *shop.Extension) {
	out.ID = 0                                // zero value
	out.UserID = arg.UserID                   // simple assign
	out.AccountID = arg.AccountID             // simple assign
	out.ExtensionNumber = arg.ExtensionNumber // simple assign
	out.ConnectionID = arg.ConnectionID       // simple assign
	out.CreatedAt = time.Time{}               // zero value
	out.UpdatedAt = time.Time{}               // zero value
}

func Convert_etelecom_Extension_shop_Extension(arg *etelecom.Extension, out *shop.Extension) *shop.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shop.Extension{}
	}
	convert_etelecom_Extension_shop_Extension(arg, out)
	return out
}

func convert_etelecom_Extension_shop_Extension(arg *etelecom.Extension, out *shop.Extension) {
	out.ID = arg.ID                           // simple assign
	out.UserID = arg.UserID                   // simple assign
	out.AccountID = arg.AccountID             // simple assign
	out.ExtensionNumber = arg.ExtensionNumber // simple assign
	out.ConnectionID = arg.ConnectionID       // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
}

func Convert_etelecom_Extensions_shop_Extensions(args []*etelecom.Extension) (outs []*shop.Extension) {
	if args == nil {
		return nil
	}
	tmps := make([]shop.Extension, len(args))
	outs = make([]*shop.Extension, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_Extension_shop_Extension(args[i], &tmps[i])
	}
	return outs
}

func Convert_shop_Extension_etelecom_Extension(arg *shop.Extension, out *etelecom.Extension) *etelecom.Extension {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.Extension{}
	}
	convert_shop_Extension_etelecom_Extension(arg, out)
	return out
}

func convert_shop_Extension_etelecom_Extension(arg *shop.Extension, out *etelecom.Extension) {
	out.ID = arg.ID                           // simple assign
	out.UserID = arg.UserID                   // simple assign
	out.AccountID = arg.AccountID             // simple assign
	out.HotlineID = 0                         // zero value
	out.ExtensionNumber = arg.ExtensionNumber // simple assign
	out.ExtensionPassword = ""                // zero value
	out.ExternalData = nil                    // zero value
	out.ConnectionID = arg.ConnectionID       // simple assign
	out.ConnectionMethod = 0                  // zero value
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.DeletedAt = time.Time{}               // zero value
}

func Convert_shop_Extensions_etelecom_Extensions(args []*shop.Extension) (outs []*etelecom.Extension) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.Extension, len(args))
	outs = make([]*etelecom.Extension, len(args))
	for i := range tmps {
		outs[i] = Convert_shop_Extension_etelecom_Extension(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/top/int/shop.ExtensionExternalData --//

func Convert_etelecom_ExtensionExternalData_shop_ExtensionExternalData(arg *etelecom.ExtensionExternalData, out *shop.ExtensionExternalData) *shop.ExtensionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shop.ExtensionExternalData{}
	}
	convert_etelecom_ExtensionExternalData_shop_ExtensionExternalData(arg, out)
	return out
}

func convert_etelecom_ExtensionExternalData_shop_ExtensionExternalData(arg *etelecom.ExtensionExternalData, out *shop.ExtensionExternalData) {
	out.ID = 0 // types do not match
}

func Convert_etelecom_ExtensionExternalDatas_shop_ExtensionExternalDatas(args []*etelecom.ExtensionExternalData) (outs []*shop.ExtensionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]shop.ExtensionExternalData, len(args))
	outs = make([]*shop.ExtensionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_etelecom_ExtensionExternalData_shop_ExtensionExternalData(args[i], &tmps[i])
	}
	return outs
}

func Convert_shop_ExtensionExternalData_etelecom_ExtensionExternalData(arg *shop.ExtensionExternalData, out *etelecom.ExtensionExternalData) *etelecom.ExtensionExternalData {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &etelecom.ExtensionExternalData{}
	}
	convert_shop_ExtensionExternalData_etelecom_ExtensionExternalData(arg, out)
	return out
}

func convert_shop_ExtensionExternalData_etelecom_ExtensionExternalData(arg *shop.ExtensionExternalData, out *etelecom.ExtensionExternalData) {
	out.ID = "" // types do not match
}

func Convert_shop_ExtensionExternalDatas_etelecom_ExtensionExternalDatas(args []*shop.ExtensionExternalData) (outs []*etelecom.ExtensionExternalData) {
	if args == nil {
		return nil
	}
	tmps := make([]etelecom.ExtensionExternalData, len(args))
	outs = make([]*etelecom.ExtensionExternalData, len(args))
	for i := range tmps {
		outs[i] = Convert_shop_ExtensionExternalData_etelecom_ExtensionExternalData(args[i], &tmps[i])
	}
	return outs
}
