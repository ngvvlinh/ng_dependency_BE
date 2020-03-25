// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	shopcustomergroupcustomermodel "etop.vn/backend/zexp/etl/main/shopcustomergroupcustomer/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer)(nil), (*customeringmodel.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		Convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg.(*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer), out.(*customeringmodel.ShopCustomerGroupCustomer))
		return nil
	})
	s.Register(([]*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer)(nil), (*[]*customeringmodel.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomers_customeringmodel_ShopCustomerGroupCustomers(arg.([]*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer))
		*out.(*[]*customeringmodel.ShopCustomerGroupCustomer) = out0
		return nil
	})
	s.Register((*customeringmodel.ShopCustomerGroupCustomer)(nil), (*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		Convert_customeringmodel_ShopCustomerGroupCustomer_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer(arg.(*customeringmodel.ShopCustomerGroupCustomer), out.(*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer))
		return nil
	})
	s.Register(([]*customeringmodel.ShopCustomerGroupCustomer)(nil), (*[]*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer)(nil), func(arg, out interface{}) error {
		out0 := Convert_customeringmodel_ShopCustomerGroupCustomers_shopcustomergroupcustomermodel_ShopCustomerGroupCustomers(arg.([]*customeringmodel.ShopCustomerGroupCustomer))
		*out.(*[]*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/shopping/customering/model.ShopCustomerGroupCustomer --//

func Convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg *shopcustomergroupcustomermodel.ShopCustomerGroupCustomer, out *customeringmodel.ShopCustomerGroupCustomer) *customeringmodel.ShopCustomerGroupCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &customeringmodel.ShopCustomerGroupCustomer{}
	}
	convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg, out)
	return out
}

func convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(arg *shopcustomergroupcustomermodel.ShopCustomerGroupCustomer, out *customeringmodel.ShopCustomerGroupCustomer) {
	out.GroupID = arg.GroupID       // simple assign
	out.CustomerID = arg.CustomerID // simple assign
	out.CreatedAt = arg.CreatedAt   // simple assign
	out.UpdatedAt = arg.UpdatedAt   // simple assign
	out.Rid = arg.Rid               // simple assign
}

func Convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomers_customeringmodel_ShopCustomerGroupCustomers(args []*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer) (outs []*customeringmodel.ShopCustomerGroupCustomer) {
	tmps := make([]customeringmodel.ShopCustomerGroupCustomer, len(args))
	outs = make([]*customeringmodel.ShopCustomerGroupCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer_customeringmodel_ShopCustomerGroupCustomer(args[i], &tmps[i])
	}
	return outs
}

func Convert_customeringmodel_ShopCustomerGroupCustomer_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer(arg *customeringmodel.ShopCustomerGroupCustomer, out *shopcustomergroupcustomermodel.ShopCustomerGroupCustomer) *shopcustomergroupcustomermodel.ShopCustomerGroupCustomer {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopcustomergroupcustomermodel.ShopCustomerGroupCustomer{}
	}
	convert_customeringmodel_ShopCustomerGroupCustomer_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer(arg, out)
	return out
}

func convert_customeringmodel_ShopCustomerGroupCustomer_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer(arg *customeringmodel.ShopCustomerGroupCustomer, out *shopcustomergroupcustomermodel.ShopCustomerGroupCustomer) {
	out.GroupID = arg.GroupID       // simple assign
	out.CustomerID = arg.CustomerID // simple assign
	out.CreatedAt = arg.CreatedAt   // simple assign
	out.UpdatedAt = arg.UpdatedAt   // simple assign
	out.Rid = arg.Rid               // simple assign
}

func Convert_customeringmodel_ShopCustomerGroupCustomers_shopcustomergroupcustomermodel_ShopCustomerGroupCustomers(args []*customeringmodel.ShopCustomerGroupCustomer) (outs []*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer) {
	tmps := make([]shopcustomergroupcustomermodel.ShopCustomerGroupCustomer, len(args))
	outs = make([]*shopcustomergroupcustomermodel.ShopCustomerGroupCustomer, len(args))
	for i := range tmps {
		outs[i] = Convert_customeringmodel_ShopCustomerGroupCustomer_shopcustomergroupcustomermodel_ShopCustomerGroupCustomer(args[i], &tmps[i])
	}
	return outs
}
