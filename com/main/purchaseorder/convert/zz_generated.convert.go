// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	purchaseorder "o.o/api/main/purchaseorder"
	catalogconvert "o.o/backend/com/main/catalog/convert"
	identityconvert "o.o/backend/com/main/identity/convert"
	purchaseordermodel "o.o/backend/com/main/purchaseorder/model"
	conversion "o.o/backend/pkg/common/conversion"
)

/*
Custom conversions:
    createPurchaseOrder    // in use
    purchaseOrder          // in use
    purchaseOrderDB        // in use
    updatePurchaseOrder    // in use

Ignored functions:
    GenerateCode     // params are not pointer to named types
    ParseCodeNorm    // not recognized
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*purchaseordermodel.PurchaseOrder)(nil), (*purchaseorder.PurchaseOrder)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel_PurchaseOrder_purchaseorder_PurchaseOrder(arg.(*purchaseordermodel.PurchaseOrder), out.(*purchaseorder.PurchaseOrder))
		return nil
	})
	s.Register(([]*purchaseordermodel.PurchaseOrder)(nil), (*[]*purchaseorder.PurchaseOrder)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel_PurchaseOrders_purchaseorder_PurchaseOrders(arg.([]*purchaseordermodel.PurchaseOrder))
		*out.(*[]*purchaseorder.PurchaseOrder) = out0
		return nil
	})
	s.Register((*purchaseorder.PurchaseOrder)(nil), (*purchaseordermodel.PurchaseOrder)(nil), func(arg, out interface{}) error {
		Convert_purchaseorder_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg.(*purchaseorder.PurchaseOrder), out.(*purchaseordermodel.PurchaseOrder))
		return nil
	})
	s.Register(([]*purchaseorder.PurchaseOrder)(nil), (*[]*purchaseordermodel.PurchaseOrder)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseorder_PurchaseOrders_purchaseordermodel_PurchaseOrders(arg.([]*purchaseorder.PurchaseOrder))
		*out.(*[]*purchaseordermodel.PurchaseOrder) = out0
		return nil
	})
	s.Register((*purchaseorder.CreatePurchaseOrderArgs)(nil), (*purchaseorder.PurchaseOrder)(nil), func(arg, out interface{}) error {
		Apply_purchaseorder_CreatePurchaseOrderArgs_purchaseorder_PurchaseOrder(arg.(*purchaseorder.CreatePurchaseOrderArgs), out.(*purchaseorder.PurchaseOrder))
		return nil
	})
	s.Register((*purchaseorder.UpdatePurchaseOrderArgs)(nil), (*purchaseorder.PurchaseOrder)(nil), func(arg, out interface{}) error {
		Apply_purchaseorder_UpdatePurchaseOrderArgs_purchaseorder_PurchaseOrder(arg.(*purchaseorder.UpdatePurchaseOrderArgs), out.(*purchaseorder.PurchaseOrder))
		return nil
	})
	s.Register((*purchaseordermodel.PurchaseOrderLine)(nil), (*purchaseorder.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel_PurchaseOrderLine_purchaseorder_PurchaseOrderLine(arg.(*purchaseordermodel.PurchaseOrderLine), out.(*purchaseorder.PurchaseOrderLine))
		return nil
	})
	s.Register(([]*purchaseordermodel.PurchaseOrderLine)(nil), (*[]*purchaseorder.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel_PurchaseOrderLines_purchaseorder_PurchaseOrderLines(arg.([]*purchaseordermodel.PurchaseOrderLine))
		*out.(*[]*purchaseorder.PurchaseOrderLine) = out0
		return nil
	})
	s.Register((*purchaseorder.PurchaseOrderLine)(nil), (*purchaseordermodel.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		Convert_purchaseorder_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg.(*purchaseorder.PurchaseOrderLine), out.(*purchaseordermodel.PurchaseOrderLine))
		return nil
	})
	s.Register(([]*purchaseorder.PurchaseOrderLine)(nil), (*[]*purchaseordermodel.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseorder_PurchaseOrderLines_purchaseordermodel_PurchaseOrderLines(arg.([]*purchaseorder.PurchaseOrderLine))
		*out.(*[]*purchaseordermodel.PurchaseOrderLine) = out0
		return nil
	})
	s.Register((*purchaseordermodel.PurchaseOrderSupplier)(nil), (*purchaseorder.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseorder_PurchaseOrderSupplier(arg.(*purchaseordermodel.PurchaseOrderSupplier), out.(*purchaseorder.PurchaseOrderSupplier))
		return nil
	})
	s.Register(([]*purchaseordermodel.PurchaseOrderSupplier)(nil), (*[]*purchaseorder.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel_PurchaseOrderSuppliers_purchaseorder_PurchaseOrderSuppliers(arg.([]*purchaseordermodel.PurchaseOrderSupplier))
		*out.(*[]*purchaseorder.PurchaseOrderSupplier) = out0
		return nil
	})
	s.Register((*purchaseorder.PurchaseOrderSupplier)(nil), (*purchaseordermodel.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		Convert_purchaseorder_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg.(*purchaseorder.PurchaseOrderSupplier), out.(*purchaseordermodel.PurchaseOrderSupplier))
		return nil
	})
	s.Register(([]*purchaseorder.PurchaseOrderSupplier)(nil), (*[]*purchaseordermodel.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseorder_PurchaseOrderSuppliers_purchaseordermodel_PurchaseOrderSuppliers(arg.([]*purchaseorder.PurchaseOrderSupplier))
		*out.(*[]*purchaseordermodel.PurchaseOrderSupplier) = out0
		return nil
	})
}

//-- convert o.o/api/main/purchaseorder.PurchaseOrder --//

func Convert_purchaseordermodel_PurchaseOrder_purchaseorder_PurchaseOrder(arg *purchaseordermodel.PurchaseOrder, out *purchaseorder.PurchaseOrder) *purchaseorder.PurchaseOrder {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseorder.PurchaseOrder{}
	}
	purchaseOrder(arg, out)
	return out
}

func convert_purchaseordermodel_PurchaseOrder_purchaseorder_PurchaseOrder(arg *purchaseordermodel.PurchaseOrder, out *purchaseorder.PurchaseOrder) {
	out.ID = arg.ID                 // simple assign
	out.ShopID = arg.ShopID         // simple assign
	out.SupplierID = arg.SupplierID // simple assign
	out.Supplier = Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseorder_PurchaseOrderSupplier(arg.Supplier, nil)
	out.InventoryVoucher = nil // zero value
	out.DiscountLines = identityconvert.Convert_sharemodel_DiscountLines_inttypes_DiscountLines(arg.DiscountLines)
	out.TotalDiscount = arg.TotalDiscount // simple assign
	out.FeeLines = identityconvert.Convert_sharemodel_FeeLines_inttypes_FeeLines(arg.FeeLines)
	out.TotalFee = arg.TotalFee       // simple assign
	out.TotalAmount = arg.TotalAmount // simple assign
	out.BasketValue = arg.BasketValue // simple assign
	out.Code = arg.Code               // simple assign
	out.CodeNorm = arg.CodeNorm       // simple assign
	out.Note = arg.Note               // simple assign
	out.Status = arg.Status           // simple assign
	out.VariantIDs = arg.VariantIDs   // simple assign
	out.Lines = Convert_purchaseordermodel_PurchaseOrderLines_purchaseorder_PurchaseOrderLines(arg.Lines)
	out.PaidAmount = 0                // zero value
	out.CreatedBy = arg.CreatedBy     // simple assign
	out.CancelReason = ""             // zero value
	out.ConfirmedAt = arg.ConfirmedAt // simple assign
	out.CancelledAt = arg.CancelledAt // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
}

func Convert_purchaseordermodel_PurchaseOrders_purchaseorder_PurchaseOrders(args []*purchaseordermodel.PurchaseOrder) (outs []*purchaseorder.PurchaseOrder) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseorder.PurchaseOrder, len(args))
	outs = make([]*purchaseorder.PurchaseOrder, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel_PurchaseOrder_purchaseorder_PurchaseOrder(args[i], &tmps[i])
	}
	return outs
}

func Convert_purchaseorder_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg *purchaseorder.PurchaseOrder, out *purchaseordermodel.PurchaseOrder) *purchaseordermodel.PurchaseOrder {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel.PurchaseOrder{}
	}
	purchaseOrderDB(arg, out)
	return out
}

func convert_purchaseorder_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg *purchaseorder.PurchaseOrder, out *purchaseordermodel.PurchaseOrder) {
	out.ID = arg.ID                 // simple assign
	out.ShopID = arg.ShopID         // simple assign
	out.SupplierID = arg.SupplierID // simple assign
	out.Supplier = Convert_purchaseorder_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg.Supplier, nil)
	out.BasketValue = arg.BasketValue // simple assign
	out.DiscountLines = identityconvert.Convert_inttypes_DiscountLines_sharemodel_DiscountLines(arg.DiscountLines)
	out.TotalDiscount = arg.TotalDiscount // simple assign
	out.FeeLines = identityconvert.Convert_inttypes_FeeLines_sharemodel_FeeLines(arg.FeeLines)
	out.TotalFee = arg.TotalFee       // simple assign
	out.TotalAmount = arg.TotalAmount // simple assign
	out.Code = arg.Code               // simple assign
	out.CodeNorm = arg.CodeNorm       // simple assign
	out.Note = arg.Note               // simple assign
	out.Status = arg.Status           // simple assign
	out.VariantIDs = arg.VariantIDs   // simple assign
	out.Lines = Convert_purchaseorder_PurchaseOrderLines_purchaseordermodel_PurchaseOrderLines(arg.Lines)
	out.CreatedBy = arg.CreatedBy     // simple assign
	out.CancelledReason = ""          // zero value
	out.ConfirmedAt = arg.ConfirmedAt // simple assign
	out.CancelledAt = arg.CancelledAt // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
	out.DeletedAt = time.Time{}       // zero value
	out.SupplierFullNameNorm = ""     // zero value
	out.SupplierPhoneNorm = ""        // zero value
	out.Rid = 0                       // zero value
}

func Convert_purchaseorder_PurchaseOrders_purchaseordermodel_PurchaseOrders(args []*purchaseorder.PurchaseOrder) (outs []*purchaseordermodel.PurchaseOrder) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel.PurchaseOrder, len(args))
	outs = make([]*purchaseordermodel.PurchaseOrder, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseorder_PurchaseOrder_purchaseordermodel_PurchaseOrder(args[i], &tmps[i])
	}
	return outs
}

func Apply_purchaseorder_CreatePurchaseOrderArgs_purchaseorder_PurchaseOrder(arg *purchaseorder.CreatePurchaseOrderArgs, out *purchaseorder.PurchaseOrder) *purchaseorder.PurchaseOrder {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseorder.PurchaseOrder{}
	}
	createPurchaseOrder(arg, out)
	return out
}

func apply_purchaseorder_CreatePurchaseOrderArgs_purchaseorder_PurchaseOrder(arg *purchaseorder.CreatePurchaseOrderArgs, out *purchaseorder.PurchaseOrder) {
	out.ID = 0                            // zero value
	out.ShopID = arg.ShopID               // simple assign
	out.SupplierID = arg.SupplierID       // simple assign
	out.Supplier = nil                    // zero value
	out.InventoryVoucher = nil            // zero value
	out.DiscountLines = arg.DiscountLines // simple assign
	out.TotalDiscount = arg.TotalDiscount // simple assign
	out.FeeLines = arg.FeeLines           // simple assign
	out.TotalFee = arg.TotalFee           // simple assign
	out.TotalAmount = arg.TotalAmount     // simple assign
	out.BasketValue = arg.BasketValue     // simple assign
	out.Code = ""                         // zero value
	out.CodeNorm = 0                      // zero value
	out.Note = arg.Note                   // simple assign
	out.Status = 0                        // zero value
	out.VariantIDs = nil                  // zero value
	out.Lines = arg.Lines                 // simple assign
	out.PaidAmount = 0                    // zero value
	out.CreatedBy = arg.CreatedBy         // simple assign
	out.CancelReason = ""                 // zero value
	out.ConfirmedAt = time.Time{}         // zero value
	out.CancelledAt = time.Time{}         // zero value
	out.CreatedAt = time.Time{}           // zero value
	out.UpdatedAt = time.Time{}           // zero value
}

func Apply_purchaseorder_UpdatePurchaseOrderArgs_purchaseorder_PurchaseOrder(arg *purchaseorder.UpdatePurchaseOrderArgs, out *purchaseorder.PurchaseOrder) *purchaseorder.PurchaseOrder {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseorder.PurchaseOrder{}
	}
	updatePurchaseOrder(arg, out)
	return out
}

func apply_purchaseorder_UpdatePurchaseOrderArgs_purchaseorder_PurchaseOrder(arg *purchaseorder.UpdatePurchaseOrderArgs, out *purchaseorder.PurchaseOrder) {
	out.ID = out.ID                                                // identifier
	out.ShopID = out.ShopID                                        // identifier
	out.SupplierID = out.SupplierID                                // no change
	out.Supplier = out.Supplier                                    // no change
	out.InventoryVoucher = out.InventoryVoucher                    // no change
	out.DiscountLines = arg.DiscountLines                          // simple assign
	out.TotalDiscount = arg.TotalDiscount.Apply(out.TotalDiscount) // apply change
	out.FeeLines = arg.FeeLines                                    // simple assign
	out.TotalFee = arg.TotalFee.Apply(out.TotalFee)                // apply change
	out.TotalAmount = arg.TotalAmount.Apply(out.TotalAmount)       // apply change
	out.BasketValue = arg.BasketValue.Apply(out.BasketValue)       // apply change
	out.Code = out.Code                                            // no change
	out.CodeNorm = out.CodeNorm                                    // no change
	out.Note = arg.Note.Apply(out.Note)                            // apply change
	out.Status = out.Status                                        // no change
	out.VariantIDs = out.VariantIDs                                // no change
	out.Lines = arg.Lines                                          // simple assign
	out.PaidAmount = out.PaidAmount                                // no change
	out.CreatedBy = out.CreatedBy                                  // no change
	out.CancelReason = out.CancelReason                            // no change
	out.ConfirmedAt = out.ConfirmedAt                              // no change
	out.CancelledAt = out.CancelledAt                              // no change
	out.CreatedAt = out.CreatedAt                                  // no change
	out.UpdatedAt = out.UpdatedAt                                  // no change
}

//-- convert o.o/api/main/purchaseorder.PurchaseOrderLine --//

func Convert_purchaseordermodel_PurchaseOrderLine_purchaseorder_PurchaseOrderLine(arg *purchaseordermodel.PurchaseOrderLine, out *purchaseorder.PurchaseOrderLine) *purchaseorder.PurchaseOrderLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseorder.PurchaseOrderLine{}
	}
	convert_purchaseordermodel_PurchaseOrderLine_purchaseorder_PurchaseOrderLine(arg, out)
	return out
}

func convert_purchaseordermodel_PurchaseOrderLine_purchaseorder_PurchaseOrderLine(arg *purchaseordermodel.PurchaseOrderLine, out *purchaseorder.PurchaseOrderLine) {
	out.VariantID = arg.VariantID       // simple assign
	out.Quantity = arg.Quantity         // simple assign
	out.PaymentPrice = arg.PaymentPrice // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.ProductName = arg.ProductName   // simple assign
	out.Code = arg.Code                 // simple assign
	out.ImageUrl = arg.ImageUrl         // simple assign
	out.Attributes = catalogconvert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(arg.Attributes)
	out.Discount = 0 // zero value
}

func Convert_purchaseordermodel_PurchaseOrderLines_purchaseorder_PurchaseOrderLines(args []*purchaseordermodel.PurchaseOrderLine) (outs []*purchaseorder.PurchaseOrderLine) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseorder.PurchaseOrderLine, len(args))
	outs = make([]*purchaseorder.PurchaseOrderLine, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel_PurchaseOrderLine_purchaseorder_PurchaseOrderLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_purchaseorder_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg *purchaseorder.PurchaseOrderLine, out *purchaseordermodel.PurchaseOrderLine) *purchaseordermodel.PurchaseOrderLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel.PurchaseOrderLine{}
	}
	convert_purchaseorder_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg, out)
	return out
}

func convert_purchaseorder_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg *purchaseorder.PurchaseOrderLine, out *purchaseordermodel.PurchaseOrderLine) {
	out.ProductName = arg.ProductName   // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.VariantID = arg.VariantID       // simple assign
	out.Quantity = arg.Quantity         // simple assign
	out.PaymentPrice = arg.PaymentPrice // simple assign
	out.Code = arg.Code                 // simple assign
	out.ImageUrl = arg.ImageUrl         // simple assign
	out.Attributes = catalogconvert.Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(arg.Attributes)
}

func Convert_purchaseorder_PurchaseOrderLines_purchaseordermodel_PurchaseOrderLines(args []*purchaseorder.PurchaseOrderLine) (outs []*purchaseordermodel.PurchaseOrderLine) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel.PurchaseOrderLine, len(args))
	outs = make([]*purchaseordermodel.PurchaseOrderLine, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseorder_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/api/main/purchaseorder.PurchaseOrderSupplier --//

func Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseorder_PurchaseOrderSupplier(arg *purchaseordermodel.PurchaseOrderSupplier, out *purchaseorder.PurchaseOrderSupplier) *purchaseorder.PurchaseOrderSupplier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseorder.PurchaseOrderSupplier{}
	}
	convert_purchaseordermodel_PurchaseOrderSupplier_purchaseorder_PurchaseOrderSupplier(arg, out)
	return out
}

func convert_purchaseordermodel_PurchaseOrderSupplier_purchaseorder_PurchaseOrderSupplier(arg *purchaseordermodel.PurchaseOrderSupplier, out *purchaseorder.PurchaseOrderSupplier) {
	out.FullName = arg.FullName                     // simple assign
	out.Phone = arg.Phone                           // simple assign
	out.Email = arg.Email                           // simple assign
	out.CompanyName = arg.CompanyName               // simple assign
	out.TaxNumber = arg.TaxNumber                   // simple assign
	out.HeadquarterAddress = arg.HeadquarterAddress // simple assign
	out.Deleted = false                             // zero value
}

func Convert_purchaseordermodel_PurchaseOrderSuppliers_purchaseorder_PurchaseOrderSuppliers(args []*purchaseordermodel.PurchaseOrderSupplier) (outs []*purchaseorder.PurchaseOrderSupplier) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseorder.PurchaseOrderSupplier, len(args))
	outs = make([]*purchaseorder.PurchaseOrderSupplier, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseorder_PurchaseOrderSupplier(args[i], &tmps[i])
	}
	return outs
}

func Convert_purchaseorder_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg *purchaseorder.PurchaseOrderSupplier, out *purchaseordermodel.PurchaseOrderSupplier) *purchaseordermodel.PurchaseOrderSupplier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel.PurchaseOrderSupplier{}
	}
	convert_purchaseorder_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg, out)
	return out
}

func convert_purchaseorder_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg *purchaseorder.PurchaseOrderSupplier, out *purchaseordermodel.PurchaseOrderSupplier) {
	out.FullName = arg.FullName                     // simple assign
	out.Phone = arg.Phone                           // simple assign
	out.Email = arg.Email                           // simple assign
	out.CompanyName = arg.CompanyName               // simple assign
	out.TaxNumber = arg.TaxNumber                   // simple assign
	out.HeadquarterAddress = arg.HeadquarterAddress // simple assign
}

func Convert_purchaseorder_PurchaseOrderSuppliers_purchaseordermodel_PurchaseOrderSuppliers(args []*purchaseorder.PurchaseOrderSupplier) (outs []*purchaseordermodel.PurchaseOrderSupplier) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel.PurchaseOrderSupplier, len(args))
	outs = make([]*purchaseordermodel.PurchaseOrderSupplier, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseorder_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(args[i], &tmps[i])
	}
	return outs
}
