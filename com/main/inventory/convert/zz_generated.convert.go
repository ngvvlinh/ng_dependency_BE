// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	inventory "etop.vn/api/main/inventory"
	catalogconvert "etop.vn/backend/com/main/catalog/convert"
	inventorymodel "etop.vn/backend/com/main/inventory/model"
	conversion "etop.vn/backend/pkg/common/conversion"
)

/*
Custom conversions:
    createInventoryVoucher    // in use
    inventoryVariant          // in use
    inventoryVariantModel     // in use
    inventoryVoucher          // in use
    inventoryVoucherModel     // in use
    updateInventoryVoucher    // in use

Ignored functions:
    GenerateCode    // params are not pointer to named types
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*inventorymodel.InventoryVariant)(nil), (*inventory.InventoryVariant)(nil), func(arg, out interface{}) error {
		Convert_inventorymodel_InventoryVariant_inventory_InventoryVariant(arg.(*inventorymodel.InventoryVariant), out.(*inventory.InventoryVariant))
		return nil
	})
	s.Register(([]*inventorymodel.InventoryVariant)(nil), (*[]*inventory.InventoryVariant)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventorymodel_InventoryVariants_inventory_InventoryVariants(arg.([]*inventorymodel.InventoryVariant))
		*out.(*[]*inventory.InventoryVariant) = out0
		return nil
	})
	s.Register((*inventory.InventoryVariant)(nil), (*inventorymodel.InventoryVariant)(nil), func(arg, out interface{}) error {
		Convert_inventory_InventoryVariant_inventorymodel_InventoryVariant(arg.(*inventory.InventoryVariant), out.(*inventorymodel.InventoryVariant))
		return nil
	})
	s.Register(([]*inventory.InventoryVariant)(nil), (*[]*inventorymodel.InventoryVariant)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventory_InventoryVariants_inventorymodel_InventoryVariants(arg.([]*inventory.InventoryVariant))
		*out.(*[]*inventorymodel.InventoryVariant) = out0
		return nil
	})
	s.Register((*inventorymodel.InventoryVoucher)(nil), (*inventory.InventoryVoucher)(nil), func(arg, out interface{}) error {
		Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(arg.(*inventorymodel.InventoryVoucher), out.(*inventory.InventoryVoucher))
		return nil
	})
	s.Register(([]*inventorymodel.InventoryVoucher)(nil), (*[]*inventory.InventoryVoucher)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventorymodel_InventoryVouchers_inventory_InventoryVouchers(arg.([]*inventorymodel.InventoryVoucher))
		*out.(*[]*inventory.InventoryVoucher) = out0
		return nil
	})
	s.Register((*inventory.InventoryVoucher)(nil), (*inventorymodel.InventoryVoucher)(nil), func(arg, out interface{}) error {
		Convert_inventory_InventoryVoucher_inventorymodel_InventoryVoucher(arg.(*inventory.InventoryVoucher), out.(*inventorymodel.InventoryVoucher))
		return nil
	})
	s.Register(([]*inventory.InventoryVoucher)(nil), (*[]*inventorymodel.InventoryVoucher)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventory_InventoryVouchers_inventorymodel_InventoryVouchers(arg.([]*inventory.InventoryVoucher))
		*out.(*[]*inventorymodel.InventoryVoucher) = out0
		return nil
	})
	s.Register((*inventory.CreateInventoryVoucherArgs)(nil), (*inventory.InventoryVoucher)(nil), func(arg, out interface{}) error {
		Apply_inventory_CreateInventoryVoucherArgs_inventory_InventoryVoucher(arg.(*inventory.CreateInventoryVoucherArgs), out.(*inventory.InventoryVoucher))
		return nil
	})
	s.Register((*inventory.UpdateInventoryVoucherArgs)(nil), (*inventory.InventoryVoucher)(nil), func(arg, out interface{}) error {
		Apply_inventory_UpdateInventoryVoucherArgs_inventory_InventoryVoucher(arg.(*inventory.UpdateInventoryVoucherArgs), out.(*inventory.InventoryVoucher))
		return nil
	})
	s.Register((*inventorymodel.InventoryVoucherItem)(nil), (*inventory.InventoryVoucherItem)(nil), func(arg, out interface{}) error {
		Convert_inventorymodel_InventoryVoucherItem_inventory_InventoryVoucherItem(arg.(*inventorymodel.InventoryVoucherItem), out.(*inventory.InventoryVoucherItem))
		return nil
	})
	s.Register(([]*inventorymodel.InventoryVoucherItem)(nil), (*[]*inventory.InventoryVoucherItem)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventorymodel_InventoryVoucherItems_inventory_InventoryVoucherItems(arg.([]*inventorymodel.InventoryVoucherItem))
		*out.(*[]*inventory.InventoryVoucherItem) = out0
		return nil
	})
	s.Register((*inventory.InventoryVoucherItem)(nil), (*inventorymodel.InventoryVoucherItem)(nil), func(arg, out interface{}) error {
		Convert_inventory_InventoryVoucherItem_inventorymodel_InventoryVoucherItem(arg.(*inventory.InventoryVoucherItem), out.(*inventorymodel.InventoryVoucherItem))
		return nil
	})
	s.Register(([]*inventory.InventoryVoucherItem)(nil), (*[]*inventorymodel.InventoryVoucherItem)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventory_InventoryVoucherItems_inventorymodel_InventoryVoucherItems(arg.([]*inventory.InventoryVoucherItem))
		*out.(*[]*inventorymodel.InventoryVoucherItem) = out0
		return nil
	})
	s.Register((*inventorymodel.Trader)(nil), (*inventory.Trader)(nil), func(arg, out interface{}) error {
		Convert_inventorymodel_Trader_inventory_Trader(arg.(*inventorymodel.Trader), out.(*inventory.Trader))
		return nil
	})
	s.Register(([]*inventorymodel.Trader)(nil), (*[]*inventory.Trader)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventorymodel_Traders_inventory_Traders(arg.([]*inventorymodel.Trader))
		*out.(*[]*inventory.Trader) = out0
		return nil
	})
	s.Register((*inventory.Trader)(nil), (*inventorymodel.Trader)(nil), func(arg, out interface{}) error {
		Convert_inventory_Trader_inventorymodel_Trader(arg.(*inventory.Trader), out.(*inventorymodel.Trader))
		return nil
	})
	s.Register(([]*inventory.Trader)(nil), (*[]*inventorymodel.Trader)(nil), func(arg, out interface{}) error {
		out0 := Convert_inventory_Traders_inventorymodel_Traders(arg.([]*inventory.Trader))
		*out.(*[]*inventorymodel.Trader) = out0
		return nil
	})
}

//-- convert etop.vn/api/main/inventory.InventoryVariant --//

func Convert_inventorymodel_InventoryVariant_inventory_InventoryVariant(arg *inventorymodel.InventoryVariant, out *inventory.InventoryVariant) *inventory.InventoryVariant {
	return inventoryVariant(arg, out)
}

func convert_inventorymodel_InventoryVariant_inventory_InventoryVariant(arg *inventorymodel.InventoryVariant, out *inventory.InventoryVariant) {
	out.ShopID = arg.ShopID                 // simple assign
	out.VariantID = arg.VariantID           // simple assign
	out.QuantityOnHand = arg.QuantityOnHand // simple assign
	out.QuantityPicked = arg.QuantityPicked // simple assign
	out.CostPrice = arg.CostPrice           // simple assign
	out.QuantitySummary = 0                 // zero value
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
}

func Convert_inventorymodel_InventoryVariants_inventory_InventoryVariants(args []*inventorymodel.InventoryVariant) (outs []*inventory.InventoryVariant) {
	tmps := make([]inventory.InventoryVariant, len(args))
	outs = make([]*inventory.InventoryVariant, len(args))
	for i := range tmps {
		outs[i] = Convert_inventorymodel_InventoryVariant_inventory_InventoryVariant(args[i], &tmps[i])
	}
	return outs
}

func Convert_inventory_InventoryVariant_inventorymodel_InventoryVariant(arg *inventory.InventoryVariant, out *inventorymodel.InventoryVariant) *inventorymodel.InventoryVariant {
	return inventoryVariantModel(arg, out)
}

func convert_inventory_InventoryVariant_inventorymodel_InventoryVariant(arg *inventory.InventoryVariant, out *inventorymodel.InventoryVariant) {
	out.ShopID = arg.ShopID                 // simple assign
	out.VariantID = arg.VariantID           // simple assign
	out.QuantityOnHand = arg.QuantityOnHand // simple assign
	out.QuantityPicked = arg.QuantityPicked // simple assign
	out.CostPrice = arg.CostPrice           // simple assign
	out.CreatedAt = arg.CreatedAt           // simple assign
	out.UpdatedAt = arg.UpdatedAt           // simple assign
}

func Convert_inventory_InventoryVariants_inventorymodel_InventoryVariants(args []*inventory.InventoryVariant) (outs []*inventorymodel.InventoryVariant) {
	tmps := make([]inventorymodel.InventoryVariant, len(args))
	outs = make([]*inventorymodel.InventoryVariant, len(args))
	for i := range tmps {
		outs[i] = Convert_inventory_InventoryVariant_inventorymodel_InventoryVariant(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/inventory.InventoryVoucher --//

func Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(arg *inventorymodel.InventoryVoucher, out *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	return inventoryVoucher(arg, out)
}

func convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(arg *inventorymodel.InventoryVoucher, out *inventory.InventoryVoucher) {
	out.ID = arg.ID                                              // simple assign
	out.ShopID = arg.ShopID                                      // simple assign
	out.Title = arg.Title                                        // simple assign
	out.Code = arg.Code                                          // simple assign
	out.CodeNorm = arg.CodeNorm                                  // simple assign
	out.CreatedBy = arg.CreatedBy                                // simple assign
	out.UpdatedBy = arg.UpdatedBy                                // simple assign
	out.CreatedAt = arg.CreatedAt                                // simple assign
	out.UpdatedAt = arg.UpdatedAt                                // simple assign
	out.ConfirmedAt = arg.ConfirmedAt                            // simple assign
	out.CancelledAt = arg.CancelledAt                            // simple assign
	out.RefID = arg.RefID                                        // simple assign
	out.RefType = inventory.InventoryRefType(arg.RefType)        // simple conversion
	out.RefName = inventory.InventoryVoucherRefName(arg.RefName) // simple conversion
	out.RefCode = arg.RefCode                                    // simple assign
	out.TraderID = arg.TraderID                                  // simple assign
	out.Trader = Convert_inventorymodel_Trader_inventory_Trader(arg.Trader, nil)
	out.TotalAmount = arg.TotalAmount   // simple assign
	out.Type = arg.Type                 // simple assign
	out.CancelReason = arg.CancelReason // simple assign
	out.Note = arg.Note                 // simple assign
	out.Lines = Convert_inventorymodel_InventoryVoucherItems_inventory_InventoryVoucherItems(arg.Lines)
	out.Status = arg.Status // simple assign
}

func Convert_inventorymodel_InventoryVouchers_inventory_InventoryVouchers(args []*inventorymodel.InventoryVoucher) (outs []*inventory.InventoryVoucher) {
	tmps := make([]inventory.InventoryVoucher, len(args))
	outs = make([]*inventory.InventoryVoucher, len(args))
	for i := range tmps {
		outs[i] = Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(args[i], &tmps[i])
	}
	return outs
}

func Convert_inventory_InventoryVoucher_inventorymodel_InventoryVoucher(arg *inventory.InventoryVoucher, out *inventorymodel.InventoryVoucher) *inventorymodel.InventoryVoucher {
	return inventoryVoucherModel(arg, out)
}

func convert_inventory_InventoryVoucher_inventorymodel_InventoryVoucher(arg *inventory.InventoryVoucher, out *inventorymodel.InventoryVoucher) {
	out.ShopID = arg.ShopID       // simple assign
	out.ID = arg.ID               // simple assign
	out.CreatedBy = arg.CreatedBy // simple assign
	out.UpdatedBy = arg.UpdatedBy // simple assign
	out.Code = arg.Code           // simple assign
	out.CodeNorm = arg.CodeNorm   // simple assign
	out.Status = arg.Status       // simple assign
	out.Note = arg.Note           // simple assign
	out.TraderID = arg.TraderID   // simple assign
	out.Trader = Convert_inventory_Trader_inventorymodel_Trader(arg.Trader, nil)
	out.TotalAmount = arg.TotalAmount // simple assign
	out.Type = arg.Type               // simple assign
	out.Lines = Convert_inventory_InventoryVoucherItems_inventorymodel_InventoryVoucherItems(arg.Lines)
	out.VariantIDs = nil                // zero value
	out.RefID = arg.RefID               // simple assign
	out.RefCode = arg.RefCode           // simple assign
	out.RefType = string(arg.RefType)   // simple conversion
	out.RefName = string(arg.RefName)   // simple conversion
	out.Title = arg.Title               // simple assign
	out.CreatedAt = arg.CreatedAt       // simple assign
	out.UpdatedAt = arg.UpdatedAt       // simple assign
	out.ConfirmedAt = arg.ConfirmedAt   // simple assign
	out.CancelledAt = arg.CancelledAt   // simple assign
	out.CancelReason = arg.CancelReason // simple assign
	out.ProductIDs = nil                // zero value
}

func Convert_inventory_InventoryVouchers_inventorymodel_InventoryVouchers(args []*inventory.InventoryVoucher) (outs []*inventorymodel.InventoryVoucher) {
	tmps := make([]inventorymodel.InventoryVoucher, len(args))
	outs = make([]*inventorymodel.InventoryVoucher, len(args))
	for i := range tmps {
		outs[i] = Convert_inventory_InventoryVoucher_inventorymodel_InventoryVoucher(args[i], &tmps[i])
	}
	return outs
}

func Apply_inventory_CreateInventoryVoucherArgs_inventory_InventoryVoucher(arg *inventory.CreateInventoryVoucherArgs, out *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inventory.InventoryVoucher{}
	}
	createInventoryVoucher(arg, out)
	return out
}

func apply_inventory_CreateInventoryVoucherArgs_inventory_InventoryVoucher(arg *inventory.CreateInventoryVoucherArgs, out *inventory.InventoryVoucher) {
	out.ID = 0                        // zero value
	out.ShopID = arg.ShopID           // simple assign
	out.Title = arg.Title             // simple assign
	out.Code = ""                     // zero value
	out.CodeNorm = 0                  // zero value
	out.CreatedBy = arg.CreatedBy     // simple assign
	out.UpdatedBy = 0                 // zero value
	out.CreatedAt = time.Time{}       // zero value
	out.UpdatedAt = time.Time{}       // zero value
	out.ConfirmedAt = time.Time{}     // zero value
	out.CancelledAt = time.Time{}     // zero value
	out.RefID = arg.RefID             // simple assign
	out.RefType = arg.RefType         // simple assign
	out.RefName = arg.RefName         // simple assign
	out.RefCode = arg.RefCode         // simple assign
	out.TraderID = arg.TraderID       // simple assign
	out.Trader = nil                  // zero value
	out.TotalAmount = arg.TotalAmount // simple assign
	out.Type = arg.Type               // simple assign
	out.CancelReason = ""             // zero value
	out.Note = arg.Note               // simple assign
	out.Lines = arg.Lines             // simple assign
	out.Status = 0                    // zero value
}

func Apply_inventory_UpdateInventoryVoucherArgs_inventory_InventoryVoucher(arg *inventory.UpdateInventoryVoucherArgs, out *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	return updateInventoryVoucher(arg, out)
}

func apply_inventory_UpdateInventoryVoucherArgs_inventory_InventoryVoucher(arg *inventory.UpdateInventoryVoucherArgs, out *inventory.InventoryVoucher) {
	out.ID = arg.ID                                 // simple assign
	out.ShopID = arg.ShopID                         // simple assign
	out.Title = arg.Title.Apply(out.Title)          // apply change
	out.Code = out.Code                             // no change
	out.CodeNorm = out.CodeNorm                     // no change
	out.CreatedBy = out.CreatedBy                   // no change
	out.UpdatedBy = arg.UpdatedBy                   // simple assign
	out.CreatedAt = out.CreatedAt                   // no change
	out.UpdatedAt = out.UpdatedAt                   // no change
	out.ConfirmedAt = out.ConfirmedAt               // no change
	out.CancelledAt = out.CancelledAt               // no change
	out.RefID = out.RefID                           // no change
	out.RefType = out.RefType                       // no change
	out.RefName = out.RefName                       // no change
	out.RefCode = out.RefCode                       // no change
	out.TraderID = arg.TraderID.Apply(out.TraderID) // apply change
	out.Trader = out.Trader                         // no change
	out.TotalAmount = arg.TotalAmount               // simple assign
	out.Type = out.Type                             // no change
	out.CancelReason = out.CancelReason             // no change
	out.Note = arg.Note.Apply(out.Note)             // apply change
	out.Lines = arg.Lines                           // simple assign
	out.Status = out.Status                         // no change
}

//-- convert etop.vn/api/main/inventory.InventoryVoucherItem --//

func Convert_inventorymodel_InventoryVoucherItem_inventory_InventoryVoucherItem(arg *inventorymodel.InventoryVoucherItem, out *inventory.InventoryVoucherItem) *inventory.InventoryVoucherItem {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inventory.InventoryVoucherItem{}
	}
	convert_inventorymodel_InventoryVoucherItem_inventory_InventoryVoucherItem(arg, out)
	return out
}

func convert_inventorymodel_InventoryVoucherItem_inventory_InventoryVoucherItem(arg *inventorymodel.InventoryVoucherItem, out *inventory.InventoryVoucherItem) {
	out.ProductID = arg.ProductID     // simple assign
	out.ProductName = arg.ProductName // simple assign
	out.VariantID = arg.VariantID     // simple assign
	out.VariantName = arg.VariantName // simple assign
	out.Quantity = arg.Quantity       // simple assign
	out.Price = arg.Price             // simple assign
	out.Code = arg.Code               // simple assign
	out.ImageURL = arg.ImageURL       // simple assign
	out.Attributes = catalogconvert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(arg.Attributes)
}

func Convert_inventorymodel_InventoryVoucherItems_inventory_InventoryVoucherItems(args []*inventorymodel.InventoryVoucherItem) (outs []*inventory.InventoryVoucherItem) {
	tmps := make([]inventory.InventoryVoucherItem, len(args))
	outs = make([]*inventory.InventoryVoucherItem, len(args))
	for i := range tmps {
		outs[i] = Convert_inventorymodel_InventoryVoucherItem_inventory_InventoryVoucherItem(args[i], &tmps[i])
	}
	return outs
}

func Convert_inventory_InventoryVoucherItem_inventorymodel_InventoryVoucherItem(arg *inventory.InventoryVoucherItem, out *inventorymodel.InventoryVoucherItem) *inventorymodel.InventoryVoucherItem {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inventorymodel.InventoryVoucherItem{}
	}
	convert_inventory_InventoryVoucherItem_inventorymodel_InventoryVoucherItem(arg, out)
	return out
}

func convert_inventory_InventoryVoucherItem_inventorymodel_InventoryVoucherItem(arg *inventory.InventoryVoucherItem, out *inventorymodel.InventoryVoucherItem) {
	out.ProductName = arg.ProductName // simple assign
	out.ProductID = arg.ProductID     // simple assign
	out.VariantID = arg.VariantID     // simple assign
	out.VariantName = arg.VariantName // simple assign
	out.Price = arg.Price             // simple assign
	out.Quantity = arg.Quantity       // simple assign
	out.Code = arg.Code               // simple assign
	out.ImageURL = arg.ImageURL       // simple assign
	out.Attributes = catalogconvert.Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(arg.Attributes)
}

func Convert_inventory_InventoryVoucherItems_inventorymodel_InventoryVoucherItems(args []*inventory.InventoryVoucherItem) (outs []*inventorymodel.InventoryVoucherItem) {
	tmps := make([]inventorymodel.InventoryVoucherItem, len(args))
	outs = make([]*inventorymodel.InventoryVoucherItem, len(args))
	for i := range tmps {
		outs[i] = Convert_inventory_InventoryVoucherItem_inventorymodel_InventoryVoucherItem(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/inventory.Trader --//

func Convert_inventorymodel_Trader_inventory_Trader(arg *inventorymodel.Trader, out *inventory.Trader) *inventory.Trader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inventory.Trader{}
	}
	convert_inventorymodel_Trader_inventory_Trader(arg, out)
	return out
}

func convert_inventorymodel_Trader_inventory_Trader(arg *inventorymodel.Trader, out *inventory.Trader) {
	out.ID = arg.ID             // simple assign
	out.Type = arg.Type         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Phone = arg.Phone       // simple assign
}

func Convert_inventorymodel_Traders_inventory_Traders(args []*inventorymodel.Trader) (outs []*inventory.Trader) {
	tmps := make([]inventory.Trader, len(args))
	outs = make([]*inventory.Trader, len(args))
	for i := range tmps {
		outs[i] = Convert_inventorymodel_Trader_inventory_Trader(args[i], &tmps[i])
	}
	return outs
}

func Convert_inventory_Trader_inventorymodel_Trader(arg *inventory.Trader, out *inventorymodel.Trader) *inventorymodel.Trader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &inventorymodel.Trader{}
	}
	convert_inventory_Trader_inventorymodel_Trader(arg, out)
	return out
}

func convert_inventory_Trader_inventorymodel_Trader(arg *inventory.Trader, out *inventorymodel.Trader) {
	out.ID = arg.ID             // simple assign
	out.Type = arg.Type         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Phone = arg.Phone       // simple assign
}

func Convert_inventory_Traders_inventorymodel_Traders(args []*inventory.Trader) (outs []*inventorymodel.Trader) {
	tmps := make([]inventorymodel.Trader, len(args))
	outs = make([]*inventorymodel.Trader, len(args))
	for i := range tmps {
		outs[i] = Convert_inventory_Trader_inventorymodel_Trader(args[i], &tmps[i])
	}
	return outs
}
