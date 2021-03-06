// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	traderingmodel "o.o/backend/com/shopping/tradering/model"
	conversion "o.o/backend/pkg/common/conversion"
	shoptradermodel "o.o/backend/zexp/etl/main/shoptrader/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shoptradermodel.ShopTrader)(nil), (*traderingmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_shoptradermodel_ShopTrader_traderingmodel_ShopTrader(arg.(*shoptradermodel.ShopTrader), out.(*traderingmodel.ShopTrader))
		return nil
	})
	s.Register(([]*shoptradermodel.ShopTrader)(nil), (*[]*traderingmodel.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_shoptradermodel_ShopTraders_traderingmodel_ShopTraders(arg.([]*shoptradermodel.ShopTrader))
		*out.(*[]*traderingmodel.ShopTrader) = out0
		return nil
	})
	s.Register((*traderingmodel.ShopTrader)(nil), (*shoptradermodel.ShopTrader)(nil), func(arg, out interface{}) error {
		Convert_traderingmodel_ShopTrader_shoptradermodel_ShopTrader(arg.(*traderingmodel.ShopTrader), out.(*shoptradermodel.ShopTrader))
		return nil
	})
	s.Register(([]*traderingmodel.ShopTrader)(nil), (*[]*shoptradermodel.ShopTrader)(nil), func(arg, out interface{}) error {
		out0 := Convert_traderingmodel_ShopTraders_shoptradermodel_ShopTraders(arg.([]*traderingmodel.ShopTrader))
		*out.(*[]*shoptradermodel.ShopTrader) = out0
		return nil
	})
}

//-- convert o.o/backend/com/shopping/tradering/model.ShopTrader --//

func Convert_shoptradermodel_ShopTrader_traderingmodel_ShopTrader(arg *shoptradermodel.ShopTrader, out *traderingmodel.ShopTrader) *traderingmodel.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &traderingmodel.ShopTrader{}
	}
	convert_shoptradermodel_ShopTrader_traderingmodel_ShopTrader(arg, out)
	return out
}

func convert_shoptradermodel_ShopTrader_traderingmodel_ShopTrader(arg *shoptradermodel.ShopTrader, out *traderingmodel.ShopTrader) {
	out.ID = arg.ID             // simple assign
	out.ShopID = arg.ShopID     // simple assign
	out.Type = arg.Type         // simple assign
	out.DeletedAt = time.Time{} // zero value
	out.Rid = arg.Rid           // simple assign
}

func Convert_shoptradermodel_ShopTraders_traderingmodel_ShopTraders(args []*shoptradermodel.ShopTrader) (outs []*traderingmodel.ShopTrader) {
	if args == nil {
		return nil
	}
	tmps := make([]traderingmodel.ShopTrader, len(args))
	outs = make([]*traderingmodel.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_shoptradermodel_ShopTrader_traderingmodel_ShopTrader(args[i], &tmps[i])
	}
	return outs
}

func Convert_traderingmodel_ShopTrader_shoptradermodel_ShopTrader(arg *traderingmodel.ShopTrader, out *shoptradermodel.ShopTrader) *shoptradermodel.ShopTrader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shoptradermodel.ShopTrader{}
	}
	convert_traderingmodel_ShopTrader_shoptradermodel_ShopTrader(arg, out)
	return out
}

func convert_traderingmodel_ShopTrader_shoptradermodel_ShopTrader(arg *traderingmodel.ShopTrader, out *shoptradermodel.ShopTrader) {
	out.ID = arg.ID         // simple assign
	out.ShopID = arg.ShopID // simple assign
	out.Type = arg.Type     // simple assign
	out.Rid = arg.Rid       // simple assign
}

func Convert_traderingmodel_ShopTraders_shoptradermodel_ShopTraders(args []*traderingmodel.ShopTrader) (outs []*shoptradermodel.ShopTrader) {
	if args == nil {
		return nil
	}
	tmps := make([]shoptradermodel.ShopTrader, len(args))
	outs = make([]*shoptradermodel.ShopTrader, len(args))
	for i := range tmps {
		outs[i] = Convert_traderingmodel_ShopTrader_shoptradermodel_ShopTrader(args[i], &tmps[i])
	}
	return outs
}
