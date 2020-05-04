// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	purchaseordermodel "o.o/backend/com/main/purchaseorder/model"
	conversion "o.o/backend/pkg/common/conversion"
	purchaseordermodel1 "o.o/backend/zexp/etl/main/purchaseorder/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*purchaseordermodel1.PurchaseOrder)(nil), (*purchaseordermodel.PurchaseOrder)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel1_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg.(*purchaseordermodel1.PurchaseOrder), out.(*purchaseordermodel.PurchaseOrder))
		return nil
	})
	s.Register(([]*purchaseordermodel1.PurchaseOrder)(nil), (*[]*purchaseordermodel.PurchaseOrder)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel1_PurchaseOrders_purchaseordermodel_PurchaseOrders(arg.([]*purchaseordermodel1.PurchaseOrder))
		*out.(*[]*purchaseordermodel.PurchaseOrder) = out0
		return nil
	})
	s.Register((*purchaseordermodel.PurchaseOrder)(nil), (*purchaseordermodel1.PurchaseOrder)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel_PurchaseOrder_purchaseordermodel1_PurchaseOrder(arg.(*purchaseordermodel.PurchaseOrder), out.(*purchaseordermodel1.PurchaseOrder))
		return nil
	})
	s.Register(([]*purchaseordermodel.PurchaseOrder)(nil), (*[]*purchaseordermodel1.PurchaseOrder)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel_PurchaseOrders_purchaseordermodel1_PurchaseOrders(arg.([]*purchaseordermodel.PurchaseOrder))
		*out.(*[]*purchaseordermodel1.PurchaseOrder) = out0
		return nil
	})
	s.Register((*purchaseordermodel1.PurchaseOrderLine)(nil), (*purchaseordermodel.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel1_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg.(*purchaseordermodel1.PurchaseOrderLine), out.(*purchaseordermodel.PurchaseOrderLine))
		return nil
	})
	s.Register(([]*purchaseordermodel1.PurchaseOrderLine)(nil), (*[]*purchaseordermodel.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel1_PurchaseOrderLines_purchaseordermodel_PurchaseOrderLines(arg.([]*purchaseordermodel1.PurchaseOrderLine))
		*out.(*[]*purchaseordermodel.PurchaseOrderLine) = out0
		return nil
	})
	s.Register((*purchaseordermodel.PurchaseOrderLine)(nil), (*purchaseordermodel1.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel_PurchaseOrderLine_purchaseordermodel1_PurchaseOrderLine(arg.(*purchaseordermodel.PurchaseOrderLine), out.(*purchaseordermodel1.PurchaseOrderLine))
		return nil
	})
	s.Register(([]*purchaseordermodel.PurchaseOrderLine)(nil), (*[]*purchaseordermodel1.PurchaseOrderLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel_PurchaseOrderLines_purchaseordermodel1_PurchaseOrderLines(arg.([]*purchaseordermodel.PurchaseOrderLine))
		*out.(*[]*purchaseordermodel1.PurchaseOrderLine) = out0
		return nil
	})
	s.Register((*purchaseordermodel1.PurchaseOrderSupplier)(nil), (*purchaseordermodel.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel1_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg.(*purchaseordermodel1.PurchaseOrderSupplier), out.(*purchaseordermodel.PurchaseOrderSupplier))
		return nil
	})
	s.Register(([]*purchaseordermodel1.PurchaseOrderSupplier)(nil), (*[]*purchaseordermodel.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel1_PurchaseOrderSuppliers_purchaseordermodel_PurchaseOrderSuppliers(arg.([]*purchaseordermodel1.PurchaseOrderSupplier))
		*out.(*[]*purchaseordermodel.PurchaseOrderSupplier) = out0
		return nil
	})
	s.Register((*purchaseordermodel.PurchaseOrderSupplier)(nil), (*purchaseordermodel1.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseordermodel1_PurchaseOrderSupplier(arg.(*purchaseordermodel.PurchaseOrderSupplier), out.(*purchaseordermodel1.PurchaseOrderSupplier))
		return nil
	})
	s.Register(([]*purchaseordermodel.PurchaseOrderSupplier)(nil), (*[]*purchaseordermodel1.PurchaseOrderSupplier)(nil), func(arg, out interface{}) error {
		out0 := Convert_purchaseordermodel_PurchaseOrderSuppliers_purchaseordermodel1_PurchaseOrderSuppliers(arg.([]*purchaseordermodel.PurchaseOrderSupplier))
		*out.(*[]*purchaseordermodel1.PurchaseOrderSupplier) = out0
		return nil
	})
}

//-- convert o.o/backend/com/main/purchaseorder/model.PurchaseOrder --//

func Convert_purchaseordermodel1_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg *purchaseordermodel1.PurchaseOrder, out *purchaseordermodel.PurchaseOrder) *purchaseordermodel.PurchaseOrder {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel.PurchaseOrder{}
	}
	convert_purchaseordermodel1_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg, out)
	return out
}

func convert_purchaseordermodel1_PurchaseOrder_purchaseordermodel_PurchaseOrder(arg *purchaseordermodel1.PurchaseOrder, out *purchaseordermodel.PurchaseOrder) {
	out.ID = arg.ID                 // simple assign
	out.ShopID = arg.ShopID         // simple assign
	out.SupplierID = arg.SupplierID // simple assign
	out.Supplier = Convert_purchaseordermodel1_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg.Supplier, nil)
	out.BasketValue = arg.BasketValue     // simple assign
	out.DiscountLines = nil               // zero value
	out.TotalDiscount = arg.TotalDiscount // simple assign
	out.FeeLines = nil                    // zero value
	out.TotalFee = arg.TotalFee           // simple assign
	out.TotalAmount = arg.TotalAmount     // simple assign
	out.Code = arg.Code                   // simple assign
	out.CodeNorm = arg.CodeNorm           // simple assign
	out.Note = arg.Note                   // simple assign
	out.Status = arg.Status               // simple assign
	out.VariantIDs = arg.VariantIDs       // simple assign
	out.Lines = Convert_purchaseordermodel1_PurchaseOrderLines_purchaseordermodel_PurchaseOrderLines(arg.Lines)
	out.CreatedBy = arg.CreatedBy             // simple assign
	out.CancelledReason = arg.CancelledReason // simple assign
	out.ConfirmedAt = arg.ConfirmedAt         // simple assign
	out.CancelledAt = arg.CancelledAt         // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.DeletedAt = time.Time{}               // zero value
	out.SupplierFullNameNorm = ""             // zero value
	out.SupplierPhoneNorm = ""                // zero value
	out.Rid = arg.Rid                         // simple assign
}

func Convert_purchaseordermodel1_PurchaseOrders_purchaseordermodel_PurchaseOrders(args []*purchaseordermodel1.PurchaseOrder) (outs []*purchaseordermodel.PurchaseOrder) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel.PurchaseOrder, len(args))
	outs = make([]*purchaseordermodel.PurchaseOrder, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel1_PurchaseOrder_purchaseordermodel_PurchaseOrder(args[i], &tmps[i])
	}
	return outs
}

func Convert_purchaseordermodel_PurchaseOrder_purchaseordermodel1_PurchaseOrder(arg *purchaseordermodel.PurchaseOrder, out *purchaseordermodel1.PurchaseOrder) *purchaseordermodel1.PurchaseOrder {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel1.PurchaseOrder{}
	}
	convert_purchaseordermodel_PurchaseOrder_purchaseordermodel1_PurchaseOrder(arg, out)
	return out
}

func convert_purchaseordermodel_PurchaseOrder_purchaseordermodel1_PurchaseOrder(arg *purchaseordermodel.PurchaseOrder, out *purchaseordermodel1.PurchaseOrder) {
	out.ID = arg.ID                 // simple assign
	out.ShopID = arg.ShopID         // simple assign
	out.SupplierID = arg.SupplierID // simple assign
	out.Supplier = Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseordermodel1_PurchaseOrderSupplier(arg.Supplier, nil)
	out.BasketValue = arg.BasketValue     // simple assign
	out.TotalDiscount = arg.TotalDiscount // simple assign
	out.TotalFee = arg.TotalFee           // simple assign
	out.TotalAmount = arg.TotalAmount     // simple assign
	out.Code = arg.Code                   // simple assign
	out.CodeNorm = arg.CodeNorm           // simple assign
	out.Note = arg.Note                   // simple assign
	out.Status = arg.Status               // simple assign
	out.VariantIDs = arg.VariantIDs       // simple assign
	out.Lines = Convert_purchaseordermodel_PurchaseOrderLines_purchaseordermodel1_PurchaseOrderLines(arg.Lines)
	out.CreatedBy = arg.CreatedBy             // simple assign
	out.CancelledReason = arg.CancelledReason // simple assign
	out.ConfirmedAt = arg.ConfirmedAt         // simple assign
	out.CancelledAt = arg.CancelledAt         // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.Rid = arg.Rid                         // simple assign
}

func Convert_purchaseordermodel_PurchaseOrders_purchaseordermodel1_PurchaseOrders(args []*purchaseordermodel.PurchaseOrder) (outs []*purchaseordermodel1.PurchaseOrder) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel1.PurchaseOrder, len(args))
	outs = make([]*purchaseordermodel1.PurchaseOrder, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel_PurchaseOrder_purchaseordermodel1_PurchaseOrder(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/purchaseorder/model.PurchaseOrderLine --//

func Convert_purchaseordermodel1_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg *purchaseordermodel1.PurchaseOrderLine, out *purchaseordermodel.PurchaseOrderLine) *purchaseordermodel.PurchaseOrderLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel.PurchaseOrderLine{}
	}
	convert_purchaseordermodel1_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg, out)
	return out
}

func convert_purchaseordermodel1_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(arg *purchaseordermodel1.PurchaseOrderLine, out *purchaseordermodel.PurchaseOrderLine) {
	out.ProductName = arg.ProductName   // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.VariantID = arg.VariantID       // simple assign
	out.Quantity = arg.Quantity         // simple assign
	out.PaymentPrice = arg.PaymentPrice // simple assign
	out.Code = arg.Code                 // simple assign
	out.ImageUrl = arg.ImageUrl         // simple assign
	out.Attributes = arg.Attributes     // simple assign
}

func Convert_purchaseordermodel1_PurchaseOrderLines_purchaseordermodel_PurchaseOrderLines(args []*purchaseordermodel1.PurchaseOrderLine) (outs []*purchaseordermodel.PurchaseOrderLine) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel.PurchaseOrderLine, len(args))
	outs = make([]*purchaseordermodel.PurchaseOrderLine, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel1_PurchaseOrderLine_purchaseordermodel_PurchaseOrderLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_purchaseordermodel_PurchaseOrderLine_purchaseordermodel1_PurchaseOrderLine(arg *purchaseordermodel.PurchaseOrderLine, out *purchaseordermodel1.PurchaseOrderLine) *purchaseordermodel1.PurchaseOrderLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel1.PurchaseOrderLine{}
	}
	convert_purchaseordermodel_PurchaseOrderLine_purchaseordermodel1_PurchaseOrderLine(arg, out)
	return out
}

func convert_purchaseordermodel_PurchaseOrderLine_purchaseordermodel1_PurchaseOrderLine(arg *purchaseordermodel.PurchaseOrderLine, out *purchaseordermodel1.PurchaseOrderLine) {
	out.ProductName = arg.ProductName   // simple assign
	out.ProductID = arg.ProductID       // simple assign
	out.VariantID = arg.VariantID       // simple assign
	out.Quantity = arg.Quantity         // simple assign
	out.PaymentPrice = arg.PaymentPrice // simple assign
	out.Code = arg.Code                 // simple assign
	out.ImageUrl = arg.ImageUrl         // simple assign
	out.Attributes = arg.Attributes     // simple assign
}

func Convert_purchaseordermodel_PurchaseOrderLines_purchaseordermodel1_PurchaseOrderLines(args []*purchaseordermodel.PurchaseOrderLine) (outs []*purchaseordermodel1.PurchaseOrderLine) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel1.PurchaseOrderLine, len(args))
	outs = make([]*purchaseordermodel1.PurchaseOrderLine, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel_PurchaseOrderLine_purchaseordermodel1_PurchaseOrderLine(args[i], &tmps[i])
	}
	return outs
}

//-- convert o.o/backend/com/main/purchaseorder/model.PurchaseOrderSupplier --//

func Convert_purchaseordermodel1_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg *purchaseordermodel1.PurchaseOrderSupplier, out *purchaseordermodel.PurchaseOrderSupplier) *purchaseordermodel.PurchaseOrderSupplier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel.PurchaseOrderSupplier{}
	}
	convert_purchaseordermodel1_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg, out)
	return out
}

func convert_purchaseordermodel1_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(arg *purchaseordermodel1.PurchaseOrderSupplier, out *purchaseordermodel.PurchaseOrderSupplier) {
	out.FullName = arg.FullName                     // simple assign
	out.Phone = arg.Phone                           // simple assign
	out.Email = arg.Email                           // simple assign
	out.CompanyName = arg.CompanyName               // simple assign
	out.TaxNumber = arg.TaxNumber                   // simple assign
	out.HeadquarterAddress = arg.HeadquarterAddress // simple assign
}

func Convert_purchaseordermodel1_PurchaseOrderSuppliers_purchaseordermodel_PurchaseOrderSuppliers(args []*purchaseordermodel1.PurchaseOrderSupplier) (outs []*purchaseordermodel.PurchaseOrderSupplier) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel.PurchaseOrderSupplier, len(args))
	outs = make([]*purchaseordermodel.PurchaseOrderSupplier, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel1_PurchaseOrderSupplier_purchaseordermodel_PurchaseOrderSupplier(args[i], &tmps[i])
	}
	return outs
}

func Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseordermodel1_PurchaseOrderSupplier(arg *purchaseordermodel.PurchaseOrderSupplier, out *purchaseordermodel1.PurchaseOrderSupplier) *purchaseordermodel1.PurchaseOrderSupplier {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &purchaseordermodel1.PurchaseOrderSupplier{}
	}
	convert_purchaseordermodel_PurchaseOrderSupplier_purchaseordermodel1_PurchaseOrderSupplier(arg, out)
	return out
}

func convert_purchaseordermodel_PurchaseOrderSupplier_purchaseordermodel1_PurchaseOrderSupplier(arg *purchaseordermodel.PurchaseOrderSupplier, out *purchaseordermodel1.PurchaseOrderSupplier) {
	out.FullName = arg.FullName                     // simple assign
	out.Phone = arg.Phone                           // simple assign
	out.Email = arg.Email                           // simple assign
	out.CompanyName = arg.CompanyName               // simple assign
	out.TaxNumber = arg.TaxNumber                   // simple assign
	out.HeadquarterAddress = arg.HeadquarterAddress // simple assign
}

func Convert_purchaseordermodel_PurchaseOrderSuppliers_purchaseordermodel1_PurchaseOrderSuppliers(args []*purchaseordermodel.PurchaseOrderSupplier) (outs []*purchaseordermodel1.PurchaseOrderSupplier) {
	if args == nil {
		return nil
	}
	tmps := make([]purchaseordermodel1.PurchaseOrderSupplier, len(args))
	outs = make([]*purchaseordermodel1.PurchaseOrderSupplier, len(args))
	for i := range tmps {
		outs[i] = Convert_purchaseordermodel_PurchaseOrderSupplier_purchaseordermodel1_PurchaseOrderSupplier(args[i], &tmps[i])
	}
	return outs
}