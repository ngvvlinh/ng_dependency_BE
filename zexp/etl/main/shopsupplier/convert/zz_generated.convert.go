// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	supplieringmodel "o.o/backend/com/shopping/suppliering/model"
	conversion "o.o/backend/pkg/common/conversion"
	shopsuppliermodel "o.o/backend/zexp/etl/main/shopsupplier/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopsuppliermodel.ShopSupplier)(nil), (*supplieringmodel.ShopSupplier)(nil), func(arg, out interface{}) error {
		Convert_shopsuppliermodel_ShopSupplier_supplieringmodel_ShopSupplier(arg.(*shopsuppliermodel.ShopSupplier), out.(*supplieringmodel.ShopSupplier))
		return nil
	})
	s.Register(([]*shopsuppliermodel.ShopSupplier)(nil), (*[]*supplieringmodel.ShopSupplier)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopsuppliermodel_ShopSuppliers_supplieringmodel_ShopSuppliers(arg.([]*shopsuppliermodel.ShopSupplier))
		*out.(*[]*supplieringmodel.ShopSupplier) = out0
		return nil
	})
	s.Register((*supplieringmodel.ShopSupplier)(nil), (*shopsuppliermodel.ShopSupplier)(nil), func(arg, out interface{}) error {
		Convert_supplieringmodel_ShopSupplier_shopsuppliermodel_ShopSupplier(arg.(*supplieringmodel.ShopSupplier), out.(*shopsuppliermodel.ShopSupplier))
		return nil
	})
	s.Register(([]*supplieringmodel.ShopSupplier)(nil), (*[]*shopsuppliermodel.ShopSupplier)(nil), func(arg, out interface{}) error {
		out0 := Convert_supplieringmodel_ShopSuppliers_shopsuppliermodel_ShopSuppliers(arg.([]*supplieringmodel.ShopSupplier))
		*out.(*[]*shopsuppliermodel.ShopSupplier) = out0
		return nil
	})
}

//-- convert o.o/backend/com/shopping/suppliering/model.ShopSupplier --//

func Convert_shopsuppliermodel_ShopSupplier_supplieringmodel_ShopSupplier(arg *shopsuppliermodel.ShopSupplier, out *supplieringmodel.ShopSupplier) *supplieringmodel.ShopSupplier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &supplieringmodel.ShopSupplier{}
	}
	convert_shopsuppliermodel_ShopSupplier_supplieringmodel_ShopSupplier(arg, out)
	return out
}

func convert_shopsuppliermodel_ShopSupplier_supplieringmodel_ShopSupplier(arg *shopsuppliermodel.ShopSupplier, out *supplieringmodel.ShopSupplier) {
	out.ID = arg.ID                               // simple assign
	out.ShopID = arg.ShopID                       // simple assign
	out.FullName = arg.FullName                   // simple assign
	out.Phone = arg.Phone                         // simple assign
	out.Email = arg.Email                         // simple assign
	out.Code = arg.Code                           // simple assign
	out.CodeNorm = 0                              // zero value
	out.CompanyName = arg.CompanyName             // simple assign
	out.CompanyNameNorm = ""                      // zero value
	out.TaxNumber = arg.TaxNumber                 // simple assign
	out.HeadquaterAddress = arg.HeadquaterAddress // simple assign
	out.Note = arg.Note                           // simple assign
	out.FullNameNorm = ""                         // zero value
	out.PhoneNorm = ""                            // zero value
	out.Status = arg.Status                       // simple assign
	out.CreatedAt = arg.CreatedAt                 // simple assign
	out.UpdatedAt = arg.UpdatedAt                 // simple assign
	out.DeletedAt = time.Time{}                   // zero value
	out.Rid = arg.Rid                             // simple assign
}

func Convert_shopsuppliermodel_ShopSuppliers_supplieringmodel_ShopSuppliers(args []*shopsuppliermodel.ShopSupplier) (outs []*supplieringmodel.ShopSupplier) {
	if args == nil {
		return nil
	}
	tmps := make([]supplieringmodel.ShopSupplier, len(args))
	outs = make([]*supplieringmodel.ShopSupplier, len(args))
	for i := range tmps {
		outs[i] = Convert_shopsuppliermodel_ShopSupplier_supplieringmodel_ShopSupplier(args[i], &tmps[i])
	}
	return outs
}

func Convert_supplieringmodel_ShopSupplier_shopsuppliermodel_ShopSupplier(arg *supplieringmodel.ShopSupplier, out *shopsuppliermodel.ShopSupplier) *shopsuppliermodel.ShopSupplier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopsuppliermodel.ShopSupplier{}
	}
	convert_supplieringmodel_ShopSupplier_shopsuppliermodel_ShopSupplier(arg, out)
	return out
}

func convert_supplieringmodel_ShopSupplier_shopsuppliermodel_ShopSupplier(arg *supplieringmodel.ShopSupplier, out *shopsuppliermodel.ShopSupplier) {
	out.ID = arg.ID                               // simple assign
	out.ShopID = arg.ShopID                       // simple assign
	out.FullName = arg.FullName                   // simple assign
	out.Phone = arg.Phone                         // simple assign
	out.Email = arg.Email                         // simple assign
	out.Code = arg.Code                           // simple assign
	out.CompanyName = arg.CompanyName             // simple assign
	out.TaxNumber = arg.TaxNumber                 // simple assign
	out.HeadquaterAddress = arg.HeadquaterAddress // simple assign
	out.Note = arg.Note                           // simple assign
	out.Status = arg.Status                       // simple assign
	out.CreatedAt = arg.CreatedAt                 // simple assign
	out.UpdatedAt = arg.UpdatedAt                 // simple assign
	out.Rid = arg.Rid                             // simple assign
}

func Convert_supplieringmodel_ShopSuppliers_shopsuppliermodel_ShopSuppliers(args []*supplieringmodel.ShopSupplier) (outs []*shopsuppliermodel.ShopSupplier) {
	if args == nil {
		return nil
	}
	tmps := make([]shopsuppliermodel.ShopSupplier, len(args))
	outs = make([]*shopsuppliermodel.ShopSupplier, len(args))
	for i := range tmps {
		outs[i] = Convert_supplieringmodel_ShopSupplier_shopsuppliermodel_ShopSupplier(args[i], &tmps[i])
	}
	return outs
}
