// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	stocktakingmodel "o.o/backend/com/main/stocktaking/model"
	conversion "o.o/backend/pkg/common/conversion"
	shopstocktakemodel "o.o/backend/zexp/etl/main/shopstocktake/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*shopstocktakemodel.ShopStocktake)(nil), (*stocktakingmodel.ShopStocktake)(nil), func(arg, out interface{}) error {
		Convert_shopstocktakemodel_ShopStocktake_stocktakingmodel_ShopStocktake(arg.(*shopstocktakemodel.ShopStocktake), out.(*stocktakingmodel.ShopStocktake))
		return nil
	})
	s.Register(([]*shopstocktakemodel.ShopStocktake)(nil), (*[]*stocktakingmodel.ShopStocktake)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopstocktakemodel_ShopStocktakes_stocktakingmodel_ShopStocktakes(arg.([]*shopstocktakemodel.ShopStocktake))
		*out.(*[]*stocktakingmodel.ShopStocktake) = out0
		return nil
	})
	s.Register((*stocktakingmodel.ShopStocktake)(nil), (*shopstocktakemodel.ShopStocktake)(nil), func(arg, out interface{}) error {
		Convert_stocktakingmodel_ShopStocktake_shopstocktakemodel_ShopStocktake(arg.(*stocktakingmodel.ShopStocktake), out.(*shopstocktakemodel.ShopStocktake))
		return nil
	})
	s.Register(([]*stocktakingmodel.ShopStocktake)(nil), (*[]*shopstocktakemodel.ShopStocktake)(nil), func(arg, out interface{}) error {
		out0 := Convert_stocktakingmodel_ShopStocktakes_shopstocktakemodel_ShopStocktakes(arg.([]*stocktakingmodel.ShopStocktake))
		*out.(*[]*shopstocktakemodel.ShopStocktake) = out0
		return nil
	})
	s.Register((*shopstocktakemodel.StocktakeLine)(nil), (*stocktakingmodel.StocktakeLine)(nil), func(arg, out interface{}) error {
		Convert_shopstocktakemodel_StocktakeLine_stocktakingmodel_StocktakeLine(arg.(*shopstocktakemodel.StocktakeLine), out.(*stocktakingmodel.StocktakeLine))
		return nil
	})
	s.Register(([]*shopstocktakemodel.StocktakeLine)(nil), (*[]*stocktakingmodel.StocktakeLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_shopstocktakemodel_StocktakeLines_stocktakingmodel_StocktakeLines(arg.([]*shopstocktakemodel.StocktakeLine))
		*out.(*[]*stocktakingmodel.StocktakeLine) = out0
		return nil
	})
	s.Register((*stocktakingmodel.StocktakeLine)(nil), (*shopstocktakemodel.StocktakeLine)(nil), func(arg, out interface{}) error {
		Convert_stocktakingmodel_StocktakeLine_shopstocktakemodel_StocktakeLine(arg.(*stocktakingmodel.StocktakeLine), out.(*shopstocktakemodel.StocktakeLine))
		return nil
	})
	s.Register(([]*stocktakingmodel.StocktakeLine)(nil), (*[]*shopstocktakemodel.StocktakeLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_stocktakingmodel_StocktakeLines_shopstocktakemodel_StocktakeLines(arg.([]*stocktakingmodel.StocktakeLine))
		*out.(*[]*shopstocktakemodel.StocktakeLine) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/stocktaking/model.ShopStocktake --//

func Convert_shopstocktakemodel_ShopStocktake_stocktakingmodel_ShopStocktake(arg *shopstocktakemodel.ShopStocktake, out *stocktakingmodel.ShopStocktake) *stocktakingmodel.ShopStocktake {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &stocktakingmodel.ShopStocktake{}
	}
	convert_shopstocktakemodel_ShopStocktake_stocktakingmodel_ShopStocktake(arg, out)
	return out
}

func convert_shopstocktakemodel_ShopStocktake_stocktakingmodel_ShopStocktake(arg *shopstocktakemodel.ShopStocktake, out *stocktakingmodel.ShopStocktake) {
	out.ID = arg.ID                       // simple assign
	out.ShopID = arg.ShopID               // simple assign
	out.TotalQuantity = arg.TotalQuantity // simple assign
	out.CreatedBy = arg.CreatedBy         // simple assign
	out.UpdatedBy = arg.UpdatedBy         // simple assign
	out.CancelReason = arg.CancelReason   // simple assign
	out.Type = arg.Type                   // simple assign
	out.Code = arg.Code                   // simple assign
	out.CodeNorm = 0                      // zero value
	out.Status = arg.Status               // simple assign
	out.CreatedAt = arg.CreatedAt         // simple assign
	out.UpdatedAt = arg.UpdatedAt         // simple assign
	out.ConfirmedAt = arg.ConfirmedAt     // simple assign
	out.CancelledAt = arg.CancelledAt     // simple assign
	out.Lines = Convert_shopstocktakemodel_StocktakeLines_stocktakingmodel_StocktakeLines(arg.Lines)
	out.Note = arg.Note             // simple assign
	out.ProductIDs = arg.ProductIDs // simple assign
	out.Rid = arg.Rid               // simple assign
}

func Convert_shopstocktakemodel_ShopStocktakes_stocktakingmodel_ShopStocktakes(args []*shopstocktakemodel.ShopStocktake) (outs []*stocktakingmodel.ShopStocktake) {
	if args == nil {
		return nil
	}
	tmps := make([]stocktakingmodel.ShopStocktake, len(args))
	outs = make([]*stocktakingmodel.ShopStocktake, len(args))
	for i := range tmps {
		outs[i] = Convert_shopstocktakemodel_ShopStocktake_stocktakingmodel_ShopStocktake(args[i], &tmps[i])
	}
	return outs
}

func Convert_stocktakingmodel_ShopStocktake_shopstocktakemodel_ShopStocktake(arg *stocktakingmodel.ShopStocktake, out *shopstocktakemodel.ShopStocktake) *shopstocktakemodel.ShopStocktake {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopstocktakemodel.ShopStocktake{}
	}
	convert_stocktakingmodel_ShopStocktake_shopstocktakemodel_ShopStocktake(arg, out)
	return out
}

func convert_stocktakingmodel_ShopStocktake_shopstocktakemodel_ShopStocktake(arg *stocktakingmodel.ShopStocktake, out *shopstocktakemodel.ShopStocktake) {
	out.ID = arg.ID                       // simple assign
	out.ShopID = arg.ShopID               // simple assign
	out.TotalQuantity = arg.TotalQuantity // simple assign
	out.CreatedBy = arg.CreatedBy         // simple assign
	out.UpdatedBy = arg.UpdatedBy         // simple assign
	out.CancelReason = arg.CancelReason   // simple assign
	out.Type = arg.Type                   // simple assign
	out.Code = arg.Code                   // simple assign
	out.Status = arg.Status               // simple assign
	out.CreatedAt = arg.CreatedAt         // simple assign
	out.UpdatedAt = arg.UpdatedAt         // simple assign
	out.ConfirmedAt = arg.ConfirmedAt     // simple assign
	out.CancelledAt = arg.CancelledAt     // simple assign
	out.Lines = Convert_stocktakingmodel_StocktakeLines_shopstocktakemodel_StocktakeLines(arg.Lines)
	out.Note = arg.Note             // simple assign
	out.ProductIDs = arg.ProductIDs // simple assign
	out.Rid = arg.Rid               // simple assign
}

func Convert_stocktakingmodel_ShopStocktakes_shopstocktakemodel_ShopStocktakes(args []*stocktakingmodel.ShopStocktake) (outs []*shopstocktakemodel.ShopStocktake) {
	if args == nil {
		return nil
	}
	tmps := make([]shopstocktakemodel.ShopStocktake, len(args))
	outs = make([]*shopstocktakemodel.ShopStocktake, len(args))
	for i := range tmps {
		outs[i] = Convert_stocktakingmodel_ShopStocktake_shopstocktakemodel_ShopStocktake(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/stocktaking/model.StocktakeLine --//

func Convert_shopstocktakemodel_StocktakeLine_stocktakingmodel_StocktakeLine(arg *shopstocktakemodel.StocktakeLine, out *stocktakingmodel.StocktakeLine) *stocktakingmodel.StocktakeLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &stocktakingmodel.StocktakeLine{}
	}
	convert_shopstocktakemodel_StocktakeLine_stocktakingmodel_StocktakeLine(arg, out)
	return out
}

func convert_shopstocktakemodel_StocktakeLine_stocktakingmodel_StocktakeLine(arg *shopstocktakemodel.StocktakeLine, out *stocktakingmodel.StocktakeLine) {
	out.ProductName = arg.ProductName // simple assign
	out.ProductID = arg.ProductID     // simple assign
	out.VariantID = arg.VariantID     // simple assign
	out.OldQuantity = arg.OldQuantity // simple assign
	out.NewQuantity = arg.NewQuantity // simple assign
	out.VariantName = arg.VariantName // simple assign
	out.Name = arg.Name               // simple assign
	out.Code = arg.Code               // simple assign
	out.ImageURL = arg.ImageURL       // simple assign
	out.Attributes = arg.Attributes   // simple assign
	out.CostPrice = arg.CostPrice     // simple assign
}

func Convert_shopstocktakemodel_StocktakeLines_stocktakingmodel_StocktakeLines(args []*shopstocktakemodel.StocktakeLine) (outs []*stocktakingmodel.StocktakeLine) {
	if args == nil {
		return nil
	}
	tmps := make([]stocktakingmodel.StocktakeLine, len(args))
	outs = make([]*stocktakingmodel.StocktakeLine, len(args))
	for i := range tmps {
		outs[i] = Convert_shopstocktakemodel_StocktakeLine_stocktakingmodel_StocktakeLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_stocktakingmodel_StocktakeLine_shopstocktakemodel_StocktakeLine(arg *stocktakingmodel.StocktakeLine, out *shopstocktakemodel.StocktakeLine) *shopstocktakemodel.StocktakeLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &shopstocktakemodel.StocktakeLine{}
	}
	convert_stocktakingmodel_StocktakeLine_shopstocktakemodel_StocktakeLine(arg, out)
	return out
}

func convert_stocktakingmodel_StocktakeLine_shopstocktakemodel_StocktakeLine(arg *stocktakingmodel.StocktakeLine, out *shopstocktakemodel.StocktakeLine) {
	out.ProductName = arg.ProductName // simple assign
	out.ProductID = arg.ProductID     // simple assign
	out.VariantID = arg.VariantID     // simple assign
	out.OldQuantity = arg.OldQuantity // simple assign
	out.NewQuantity = arg.NewQuantity // simple assign
	out.VariantName = arg.VariantName // simple assign
	out.Name = arg.Name               // simple assign
	out.Code = arg.Code               // simple assign
	out.ImageURL = arg.ImageURL       // simple assign
	out.Attributes = arg.Attributes   // simple assign
	out.CostPrice = arg.CostPrice     // simple assign
}

func Convert_stocktakingmodel_StocktakeLines_shopstocktakemodel_StocktakeLines(args []*stocktakingmodel.StocktakeLine) (outs []*shopstocktakemodel.StocktakeLine) {
	if args == nil {
		return nil
	}
	tmps := make([]shopstocktakemodel.StocktakeLine, len(args))
	outs = make([]*shopstocktakemodel.StocktakeLine, len(args))
	for i := range tmps {
		outs[i] = Convert_stocktakingmodel_StocktakeLine_shopstocktakemodel_StocktakeLine(args[i], &tmps[i])
	}
	return outs
}
