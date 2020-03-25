// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	tradering "etop.vn/api/shopping/tradering"
	traderingmodel "etop.vn/backend/com/shopping/tradering/model"
	conversion "etop.vn/backend/pkg/common/conversion"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*traderingmodel.ShopTrader)(nil), (*tradering.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_traderingmodel_ShopTrader_tradering_ShopTrader(arg.(*traderingmodel.ShopTrader), out.(*tradering.ShopTrader))
		return nil
	})
	s.Register(([]*traderingmodel.ShopTrader)(nil), (*[]*tradering.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_traderingmodel_ShopTraders_tradering_ShopTraders(arg.([]*traderingmodel.ShopTrader))
		*out.(*[]*tradering.ShopTrader) = out0
		return nil
	})
	s.Register((*tradering.ShopTrader)(nil), (*traderingmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_tradering_ShopTrader_traderingmodel_ShopTrader(arg.(*tradering.ShopTrader), out.(*traderingmodel.ShopTrader))
		return nil
	})
	s.Register(([]*tradering.ShopTrader)(nil), (*[]*traderingmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_tradering_ShopTraders_traderingmodel_ShopTraders(arg.([]*tradering.ShopTrader))
		*out.(*[]*traderingmodel.ShopTrader) = out0
		return nil
	})
}

//-- convert etop.vn/api/shopping/tradering.ShopTrader --//

func Convert_traderingmodel_ShopTrader_tradering_ShopTrader(arg *traderingmodel.ShopTrader, out *tradering.ShopTrader) *tradering.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &tradering.ShopTrader{}
	}
	convert_traderingmodel_ShopTrader_tradering_ShopTrader(arg, out)
	return out
}

func convert_traderingmodel_ShopTrader_tradering_ShopTrader(arg *traderingmodel.ShopTrader, out *tradering.ShopTrader) {
	out.ID = arg.ID         // simple assign
	out.ShopID = arg.ShopID // simple assign
	out.Type = arg.Type     // simple assign
	out.FullName = ""       // zero value
	out.Phone = ""          // zero value
}

func Convert_traderingmodel_ShopTraders_tradering_ShopTraders(args []*traderingmodel.ShopTrader) (outs []*tradering.ShopTrader) {
	tmps := make([]tradering.ShopTrader, len(args))
	outs = make([]*tradering.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_traderingmodel_ShopTrader_tradering_ShopTrader(args[i], &tmps[i])
	}
	return outs
}

func Convert_tradering_ShopTrader_traderingmodel_ShopTrader(arg *tradering.ShopTrader, out *traderingmodel.ShopTrader) *traderingmodel.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &traderingmodel.ShopTrader{}
	}
	convert_tradering_ShopTrader_traderingmodel_ShopTrader(arg, out)
	return out
}

func convert_tradering_ShopTrader_traderingmodel_ShopTrader(arg *tradering.ShopTrader, out *traderingmodel.ShopTrader) {
	out.ID = arg.ID             // simple assign
	out.ShopID = arg.ShopID     // simple assign
	out.Type = arg.Type         // simple assign
	out.DeletedAt = time.Time{} // zero value
	out.Rid = 0                 // zero value
}

func Convert_tradering_ShopTraders_traderingmodel_ShopTraders(args []*tradering.ShopTrader) (outs []*traderingmodel.ShopTrader) {
	tmps := make([]traderingmodel.ShopTrader, len(args))
	outs = make([]*traderingmodel.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_tradering_ShopTrader_traderingmodel_ShopTrader(args[i], &tmps[i])
	}
	return outs
}
