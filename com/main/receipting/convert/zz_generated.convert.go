// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	receipting "etop.vn/api/main/receipting"
	status3 "etop.vn/api/top/types/etc/status3"
	receiptingmodel "etop.vn/backend/com/main/receipting/model"
	conversion "etop.vn/backend/pkg/common/conversion"
)

/*
Custom conversions:
    createReceipt    // in use
    receiptDB        // in use
    updateReceipt    // in use

Ignored functions:
    GenerateCode     // params are not pointer to named types
    ParseCodeNorm    // not recognized
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*receiptingmodel.Receipt)(nil), (*receipting.Receipt)(nil), func(arg, out interface{}) error {
		Convert_receiptingmodel_Receipt_receipting_Receipt(arg.(*receiptingmodel.Receipt), out.(*receipting.Receipt))
		return nil
	})
	s.Register(([]*receiptingmodel.Receipt)(nil), (*[]*receipting.Receipt)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptingmodel_Receipts_receipting_Receipts(arg.([]*receiptingmodel.Receipt))
		*out.(*[]*receipting.Receipt) = out0
		return nil
	})
	s.Register((*receipting.Receipt)(nil), (*receiptingmodel.Receipt)(nil), func(arg, out interface{}) error {
		Convert_receipting_Receipt_receiptingmodel_Receipt(arg.(*receipting.Receipt), out.(*receiptingmodel.Receipt))
		return nil
	})
	s.Register(([]*receipting.Receipt)(nil), (*[]*receiptingmodel.Receipt)(nil), func(arg, out interface{}) error {
		out0 := Convert_receipting_Receipts_receiptingmodel_Receipts(arg.([]*receipting.Receipt))
		*out.(*[]*receiptingmodel.Receipt) = out0
		return nil
	})
	s.Register((*receipting.CreateReceiptArgs)(nil), (*receipting.Receipt)(nil), func(arg, out interface{}) error {
		Apply_receipting_CreateReceiptArgs_receipting_Receipt(arg.(*receipting.CreateReceiptArgs), out.(*receipting.Receipt))
		return nil
	})
	s.Register((*receipting.UpdateReceiptArgs)(nil), (*receipting.Receipt)(nil), func(arg, out interface{}) error {
		Apply_receipting_UpdateReceiptArgs_receipting_Receipt(arg.(*receipting.UpdateReceiptArgs), out.(*receipting.Receipt))
		return nil
	})
	s.Register((*receiptingmodel.ReceiptLine)(nil), (*receipting.ReceiptLine)(nil), func(arg, out interface{}) error {
		Convert_receiptingmodel_ReceiptLine_receipting_ReceiptLine(arg.(*receiptingmodel.ReceiptLine), out.(*receipting.ReceiptLine))
		return nil
	})
	s.Register(([]*receiptingmodel.ReceiptLine)(nil), (*[]*receipting.ReceiptLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptingmodel_ReceiptLines_receipting_ReceiptLines(arg.([]*receiptingmodel.ReceiptLine))
		*out.(*[]*receipting.ReceiptLine) = out0
		return nil
	})
	s.Register((*receipting.ReceiptLine)(nil), (*receiptingmodel.ReceiptLine)(nil), func(arg, out interface{}) error {
		Convert_receipting_ReceiptLine_receiptingmodel_ReceiptLine(arg.(*receipting.ReceiptLine), out.(*receiptingmodel.ReceiptLine))
		return nil
	})
	s.Register(([]*receipting.ReceiptLine)(nil), (*[]*receiptingmodel.ReceiptLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_receipting_ReceiptLines_receiptingmodel_ReceiptLines(arg.([]*receipting.ReceiptLine))
		*out.(*[]*receiptingmodel.ReceiptLine) = out0
		return nil
	})
	s.Register((*receiptingmodel.Trader)(nil), (*receipting.Trader)(nil), func(arg, out interface{}) error {
		Convert_receiptingmodel_Trader_receipting_Trader(arg.(*receiptingmodel.Trader), out.(*receipting.Trader))
		return nil
	})
	s.Register(([]*receiptingmodel.Trader)(nil), (*[]*receipting.Trader)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptingmodel_Traders_receipting_Traders(arg.([]*receiptingmodel.Trader))
		*out.(*[]*receipting.Trader) = out0
		return nil
	})
	s.Register((*receipting.Trader)(nil), (*receiptingmodel.Trader)(nil), func(arg, out interface{}) error {
		Convert_receipting_Trader_receiptingmodel_Trader(arg.(*receipting.Trader), out.(*receiptingmodel.Trader))
		return nil
	})
	s.Register(([]*receipting.Trader)(nil), (*[]*receiptingmodel.Trader)(nil), func(arg, out interface{}) error {
		out0 := Convert_receipting_Traders_receiptingmodel_Traders(arg.([]*receipting.Trader))
		*out.(*[]*receiptingmodel.Trader) = out0
		return nil
	})
}

//-- convert etop.vn/api/main/receipting.Receipt --//

func Convert_receiptingmodel_Receipt_receipting_Receipt(arg *receiptingmodel.Receipt, out *receipting.Receipt) *receipting.Receipt {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receipting.Receipt{}
	}
	convert_receiptingmodel_Receipt_receipting_Receipt(arg, out)
	return out
}

func convert_receiptingmodel_Receipt_receipting_Receipt(arg *receiptingmodel.Receipt, out *receipting.Receipt) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.TraderID = arg.TraderID       // simple assign
	out.Code = arg.Code               // simple assign
	out.CodeNorm = arg.CodeNorm       // simple assign
	out.Title = arg.Title             // simple assign
	out.Type = arg.Type               // simple assign
	out.Description = arg.Description // simple assign
	out.Amount = arg.Amount           // simple assign
	out.Status = arg.Status           // simple assign
	out.LedgerID = arg.LedgerID       // simple assign
	out.RefIDs = arg.RefIDs           // simple assign
	out.RefType = arg.RefType         // simple assign
	out.Lines = Convert_receiptingmodel_ReceiptLines_receipting_ReceiptLines(arg.Lines)
	out.Trader = Convert_receiptingmodel_Trader_receipting_Trader(arg.Trader, nil)
	out.PaidAt = arg.PaidAt           // simple assign
	out.ConfirmedAt = arg.ConfirmedAt // simple assign
	out.CancelledAt = arg.CancelledAt // simple assign
	out.CreatedBy = arg.CreatedBy     // simple assign
	out.Mode = 0                      // zero value
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
}

func Convert_receiptingmodel_Receipts_receipting_Receipts(args []*receiptingmodel.Receipt) (outs []*receipting.Receipt) {
	tmps := make([]receipting.Receipt, len(args))
	outs = make([]*receipting.Receipt, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptingmodel_Receipt_receipting_Receipt(args[i], &tmps[i])
	}
	return outs
}

func Convert_receipting_Receipt_receiptingmodel_Receipt(arg *receipting.Receipt, out *receiptingmodel.Receipt) *receiptingmodel.Receipt {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptingmodel.Receipt{}
	}
	receiptDB(arg, out)
	return out
}

func convert_receipting_Receipt_receiptingmodel_Receipt(arg *receipting.Receipt, out *receiptingmodel.Receipt) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.TraderID = arg.TraderID       // simple assign
	out.Code = arg.Code               // simple assign
	out.CodeNorm = arg.CodeNorm       // simple assign
	out.Title = arg.Title             // simple assign
	out.Type = arg.Type               // simple assign
	out.Description = arg.Description // simple assign
	out.TraderFullNameNorm = ""       // zero value
	out.TraderPhoneNorm = ""          // zero value
	out.TraderType = ""               // zero value
	out.Amount = arg.Amount           // simple assign
	out.Status = arg.Status           // simple assign
	out.RefIDs = arg.RefIDs           // simple assign
	out.RefType = arg.RefType         // simple assign
	out.Lines = Convert_receipting_ReceiptLines_receiptingmodel_ReceiptLines(arg.Lines)
	out.LedgerID = arg.LedgerID // simple assign
	out.Trader = Convert_receipting_Trader_receiptingmodel_Trader(arg.Trader, nil)
	out.CancelledReason = ""          // zero value
	out.CreatedType = ""              // zero value
	out.CreatedBy = arg.CreatedBy     // simple assign
	out.PaidAt = arg.PaidAt           // simple assign
	out.ConfirmedAt = arg.ConfirmedAt // simple assign
	out.CancelledAt = arg.CancelledAt // simple assign
	out.CreatedAt = arg.CreatedAt     // simple assign
	out.UpdatedAt = arg.UpdatedAt     // simple assign
	out.DeletedAt = time.Time{}       // zero value
}

func Convert_receipting_Receipts_receiptingmodel_Receipts(args []*receipting.Receipt) (outs []*receiptingmodel.Receipt) {
	tmps := make([]receiptingmodel.Receipt, len(args))
	outs = make([]*receiptingmodel.Receipt, len(args))
	for i := range tmps {
		outs[i] = Convert_receipting_Receipt_receiptingmodel_Receipt(args[i], &tmps[i])
	}
	return outs
}

func Apply_receipting_CreateReceiptArgs_receipting_Receipt(arg *receipting.CreateReceiptArgs, out *receipting.Receipt) *receipting.Receipt {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receipting.Receipt{}
	}
	createReceipt(arg, out)
	return out
}

func apply_receipting_CreateReceiptArgs_receipting_Receipt(arg *receipting.CreateReceiptArgs, out *receipting.Receipt) {
	out.ID = 0                              // zero value
	out.ShopID = arg.ShopID                 // simple assign
	out.TraderID = arg.TraderID             // simple assign
	out.Code = ""                           // zero value
	out.CodeNorm = 0                        // zero value
	out.Title = arg.Title                   // simple assign
	out.Type = arg.Type                     // simple assign
	out.Description = arg.Description       // simple assign
	out.Amount = arg.Amount                 // simple assign
	out.Status = status3.Status(arg.Status) // simple conversion
	out.LedgerID = arg.LedgerID             // simple assign
	out.RefIDs = arg.RefIDs                 // simple assign
	out.RefType = arg.RefType               // simple assign
	out.Lines = arg.Lines                   // simple assign
	out.Trader = arg.Trader                 // simple assign
	out.PaidAt = arg.PaidAt                 // simple assign
	out.ConfirmedAt = arg.ConfirmedAt       // simple assign
	out.CancelledAt = time.Time{}           // zero value
	out.CreatedBy = arg.CreatedBy           // simple assign
	out.Mode = arg.Mode                     // simple assign
	out.CreatedAt = time.Time{}             // zero value
	out.UpdatedAt = time.Time{}             // zero value
}

func Apply_receipting_UpdateReceiptArgs_receipting_Receipt(arg *receipting.UpdateReceiptArgs, out *receipting.Receipt) *receipting.Receipt {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receipting.Receipt{}
	}
	updateReceipt(arg, out)
	return out
}

func apply_receipting_UpdateReceiptArgs_receipting_Receipt(arg *receipting.UpdateReceiptArgs, out *receipting.Receipt) {
	out.ID = out.ID                                          // identifier
	out.ShopID = out.ShopID                                  // identifier
	out.TraderID = arg.TraderID.Apply(out.TraderID)          // apply change
	out.Code = out.Code                                      // no change
	out.CodeNorm = out.CodeNorm                              // no change
	out.Title = arg.Title.Apply(out.Title)                   // apply change
	out.Type = out.Type                                      // no change
	out.Description = arg.Description.Apply(out.Description) // apply change
	out.Amount = arg.Amount.Apply(out.Amount)                // apply change
	out.Status = out.Status                                  // no change
	out.LedgerID = arg.LedgerID.Apply(out.LedgerID)          // apply change
	out.RefIDs = arg.RefIDs                                  // simple assign
	out.RefType = arg.RefType.Apply(out.RefType)             // apply change
	out.Lines = arg.Lines                                    // simple assign
	out.Trader = arg.Trader                                  // simple assign
	out.PaidAt = arg.PaidAt                                  // simple assign
	out.ConfirmedAt = out.ConfirmedAt                        // no change
	out.CancelledAt = out.CancelledAt                        // no change
	out.CreatedBy = out.CreatedBy                            // no change
	out.Mode = out.Mode                                      // no change
	out.CreatedAt = out.CreatedAt                            // no change
	out.UpdatedAt = out.UpdatedAt                            // no change
}

//-- convert etop.vn/api/main/receipting.ReceiptLine --//

func Convert_receiptingmodel_ReceiptLine_receipting_ReceiptLine(arg *receiptingmodel.ReceiptLine, out *receipting.ReceiptLine) *receipting.ReceiptLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receipting.ReceiptLine{}
	}
	convert_receiptingmodel_ReceiptLine_receipting_ReceiptLine(arg, out)
	return out
}

func convert_receiptingmodel_ReceiptLine_receipting_ReceiptLine(arg *receiptingmodel.ReceiptLine, out *receipting.ReceiptLine) {
	out.RefID = arg.RefID   // simple assign
	out.Title = arg.Title   // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_receiptingmodel_ReceiptLines_receipting_ReceiptLines(args []*receiptingmodel.ReceiptLine) (outs []*receipting.ReceiptLine) {
	tmps := make([]receipting.ReceiptLine, len(args))
	outs = make([]*receipting.ReceiptLine, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptingmodel_ReceiptLine_receipting_ReceiptLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_receipting_ReceiptLine_receiptingmodel_ReceiptLine(arg *receipting.ReceiptLine, out *receiptingmodel.ReceiptLine) *receiptingmodel.ReceiptLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptingmodel.ReceiptLine{}
	}
	convert_receipting_ReceiptLine_receiptingmodel_ReceiptLine(arg, out)
	return out
}

func convert_receipting_ReceiptLine_receiptingmodel_ReceiptLine(arg *receipting.ReceiptLine, out *receiptingmodel.ReceiptLine) {
	out.RefID = arg.RefID   // simple assign
	out.Title = arg.Title   // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_receipting_ReceiptLines_receiptingmodel_ReceiptLines(args []*receipting.ReceiptLine) (outs []*receiptingmodel.ReceiptLine) {
	tmps := make([]receiptingmodel.ReceiptLine, len(args))
	outs = make([]*receiptingmodel.ReceiptLine, len(args))
	for i := range tmps {
		outs[i] = Convert_receipting_ReceiptLine_receiptingmodel_ReceiptLine(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/api/main/receipting.Trader --//

func Convert_receiptingmodel_Trader_receipting_Trader(arg *receiptingmodel.Trader, out *receipting.Trader) *receipting.Trader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receipting.Trader{}
	}
	convert_receiptingmodel_Trader_receipting_Trader(arg, out)
	return out
}

func convert_receiptingmodel_Trader_receipting_Trader(arg *receiptingmodel.Trader, out *receipting.Trader) {
	out.ID = arg.ID             // simple assign
	out.Type = arg.Type         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Phone = arg.Phone       // simple assign
}

func Convert_receiptingmodel_Traders_receipting_Traders(args []*receiptingmodel.Trader) (outs []*receipting.Trader) {
	tmps := make([]receipting.Trader, len(args))
	outs = make([]*receipting.Trader, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptingmodel_Trader_receipting_Trader(args[i], &tmps[i])
	}
	return outs
}

func Convert_receipting_Trader_receiptingmodel_Trader(arg *receipting.Trader, out *receiptingmodel.Trader) *receiptingmodel.Trader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptingmodel.Trader{}
	}
	convert_receipting_Trader_receiptingmodel_Trader(arg, out)
	return out
}

func convert_receipting_Trader_receiptingmodel_Trader(arg *receipting.Trader, out *receiptingmodel.Trader) {
	out.ID = arg.ID             // simple assign
	out.Type = arg.Type         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Phone = arg.Phone       // simple assign
}

func Convert_receipting_Traders_receiptingmodel_Traders(args []*receipting.Trader) (outs []*receiptingmodel.Trader) {
	tmps := make([]receiptingmodel.Trader, len(args))
	outs = make([]*receiptingmodel.Trader, len(args))
	for i := range tmps {
		outs[i] = Convert_receipting_Trader_receiptingmodel_Trader(args[i], &tmps[i])
	}
	return outs
}
