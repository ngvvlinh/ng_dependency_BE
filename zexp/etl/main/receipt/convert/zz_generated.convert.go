// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	receiptingmodel "etop.vn/backend/com/main/receipting/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	receiptmodel "etop.vn/backend/zexp/etl/main/receipt/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*receiptmodel.Receipt)(nil), (*receiptingmodel.Receipt)(nil), func(arg, out interface{}) error {
		Convert_receiptmodel_Receipt_receiptingmodel_Receipt(arg.(*receiptmodel.Receipt), out.(*receiptingmodel.Receipt))
		return nil
	})
	s.Register(([]*receiptmodel.Receipt)(nil), (*[]*receiptingmodel.Receipt)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptmodel_Receipts_receiptingmodel_Receipts(arg.([]*receiptmodel.Receipt))
		*out.(*[]*receiptingmodel.Receipt) = out0
		return nil
	})
	s.Register((*receiptingmodel.Receipt)(nil), (*receiptmodel.Receipt)(nil), func(arg, out interface{}) error {
		Convert_receiptingmodel_Receipt_receiptmodel_Receipt(arg.(*receiptingmodel.Receipt), out.(*receiptmodel.Receipt))
		return nil
	})
	s.Register(([]*receiptingmodel.Receipt)(nil), (*[]*receiptmodel.Receipt)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptingmodel_Receipts_receiptmodel_Receipts(arg.([]*receiptingmodel.Receipt))
		*out.(*[]*receiptmodel.Receipt) = out0
		return nil
	})
	s.Register((*receiptmodel.ReceiptLine)(nil), (*receiptingmodel.ReceiptLine)(nil), func(arg, out interface{}) error {
		Convert_receiptmodel_ReceiptLine_receiptingmodel_ReceiptLine(arg.(*receiptmodel.ReceiptLine), out.(*receiptingmodel.ReceiptLine))
		return nil
	})
	s.Register(([]*receiptmodel.ReceiptLine)(nil), (*[]*receiptingmodel.ReceiptLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptmodel_ReceiptLines_receiptingmodel_ReceiptLines(arg.([]*receiptmodel.ReceiptLine))
		*out.(*[]*receiptingmodel.ReceiptLine) = out0
		return nil
	})
	s.Register((*receiptingmodel.ReceiptLine)(nil), (*receiptmodel.ReceiptLine)(nil), func(arg, out interface{}) error {
		Convert_receiptingmodel_ReceiptLine_receiptmodel_ReceiptLine(arg.(*receiptingmodel.ReceiptLine), out.(*receiptmodel.ReceiptLine))
		return nil
	})
	s.Register(([]*receiptingmodel.ReceiptLine)(nil), (*[]*receiptmodel.ReceiptLine)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptingmodel_ReceiptLines_receiptmodel_ReceiptLines(arg.([]*receiptingmodel.ReceiptLine))
		*out.(*[]*receiptmodel.ReceiptLine) = out0
		return nil
	})
	s.Register((*receiptmodel.Trader)(nil), (*receiptingmodel.Trader)(nil), func(arg, out interface{}) error {
		Convert_receiptmodel_Trader_receiptingmodel_Trader(arg.(*receiptmodel.Trader), out.(*receiptingmodel.Trader))
		return nil
	})
	s.Register(([]*receiptmodel.Trader)(nil), (*[]*receiptingmodel.Trader)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptmodel_Traders_receiptingmodel_Traders(arg.([]*receiptmodel.Trader))
		*out.(*[]*receiptingmodel.Trader) = out0
		return nil
	})
	s.Register((*receiptingmodel.Trader)(nil), (*receiptmodel.Trader)(nil), func(arg, out interface{}) error {
		Convert_receiptingmodel_Trader_receiptmodel_Trader(arg.(*receiptingmodel.Trader), out.(*receiptmodel.Trader))
		return nil
	})
	s.Register(([]*receiptingmodel.Trader)(nil), (*[]*receiptmodel.Trader)(nil), func(arg, out interface{}) error {
		out0 := Convert_receiptingmodel_Traders_receiptmodel_Traders(arg.([]*receiptingmodel.Trader))
		*out.(*[]*receiptmodel.Trader) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/main/receipting/model.Receipt --//

func Convert_receiptmodel_Receipt_receiptingmodel_Receipt(arg *receiptmodel.Receipt, out *receiptingmodel.Receipt) *receiptingmodel.Receipt {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptingmodel.Receipt{}
	}
	convert_receiptmodel_Receipt_receiptingmodel_Receipt(arg, out)
	return out
}

func convert_receiptmodel_Receipt_receiptingmodel_Receipt(arg *receiptmodel.Receipt, out *receiptingmodel.Receipt) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.TraderID = arg.TraderID       // simple assign
	out.Code = arg.Code               // simple assign
	out.CodeNorm = 0                  // zero value
	out.Title = arg.Title             // simple assign
	out.Type = arg.Type               // simple assign
	out.Description = arg.Description // simple assign
	out.TraderFullNameNorm = ""       // zero value
	out.TraderPhoneNorm = ""          // zero value
	out.TraderType = arg.TraderType   // simple assign
	out.Amount = arg.Amount           // simple assign
	out.Status = arg.Status           // simple assign
	out.RefIDs = arg.RefIDs           // simple assign
	out.RefType = arg.RefType         // simple assign
	out.Lines = Convert_receiptmodel_ReceiptLines_receiptingmodel_ReceiptLines(arg.Lines)
	out.LedgerID = arg.LedgerID // simple assign
	out.Trader = Convert_receiptmodel_Trader_receiptingmodel_Trader(arg.Trader, nil)
	out.CancelledReason = arg.CancelledReason // simple assign
	out.CreatedType = arg.CreatedType         // simple assign
	out.CreatedBy = arg.CreatedBy             // simple assign
	out.PaidAt = arg.PaidAt                   // simple assign
	out.ConfirmedAt = arg.ConfirmedAt         // simple assign
	out.CancelledAt = arg.CancelledAt         // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.DeletedAt = time.Time{}               // zero value
	out.Rid = arg.Rid                         // simple assign
}

func Convert_receiptmodel_Receipts_receiptingmodel_Receipts(args []*receiptmodel.Receipt) (outs []*receiptingmodel.Receipt) {
	tmps := make([]receiptingmodel.Receipt, len(args))
	outs = make([]*receiptingmodel.Receipt, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptmodel_Receipt_receiptingmodel_Receipt(args[i], &tmps[i])
	}
	return outs
}

func Convert_receiptingmodel_Receipt_receiptmodel_Receipt(arg *receiptingmodel.Receipt, out *receiptmodel.Receipt) *receiptmodel.Receipt {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptmodel.Receipt{}
	}
	convert_receiptingmodel_Receipt_receiptmodel_Receipt(arg, out)
	return out
}

func convert_receiptingmodel_Receipt_receiptmodel_Receipt(arg *receiptingmodel.Receipt, out *receiptmodel.Receipt) {
	out.ID = arg.ID                   // simple assign
	out.ShopID = arg.ShopID           // simple assign
	out.TraderID = arg.TraderID       // simple assign
	out.Code = arg.Code               // simple assign
	out.Title = arg.Title             // simple assign
	out.Type = arg.Type               // simple assign
	out.Description = arg.Description // simple assign
	out.TraderType = arg.TraderType   // simple assign
	out.Amount = arg.Amount           // simple assign
	out.Status = arg.Status           // simple assign
	out.RefIDs = arg.RefIDs           // simple assign
	out.RefType = arg.RefType         // simple assign
	out.Lines = Convert_receiptingmodel_ReceiptLines_receiptmodel_ReceiptLines(arg.Lines)
	out.LedgerID = arg.LedgerID // simple assign
	out.Trader = Convert_receiptingmodel_Trader_receiptmodel_Trader(arg.Trader, nil)
	out.CancelledReason = arg.CancelledReason // simple assign
	out.CreatedType = arg.CreatedType         // simple assign
	out.CreatedBy = arg.CreatedBy             // simple assign
	out.PaidAt = arg.PaidAt                   // simple assign
	out.ConfirmedAt = arg.ConfirmedAt         // simple assign
	out.CancelledAt = arg.CancelledAt         // simple assign
	out.CreatedAt = arg.CreatedAt             // simple assign
	out.UpdatedAt = arg.UpdatedAt             // simple assign
	out.Rid = arg.Rid                         // simple assign
}

func Convert_receiptingmodel_Receipts_receiptmodel_Receipts(args []*receiptingmodel.Receipt) (outs []*receiptmodel.Receipt) {
	tmps := make([]receiptmodel.Receipt, len(args))
	outs = make([]*receiptmodel.Receipt, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptingmodel_Receipt_receiptmodel_Receipt(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/backend/com/main/receipting/model.ReceiptLine --//

func Convert_receiptmodel_ReceiptLine_receiptingmodel_ReceiptLine(arg *receiptmodel.ReceiptLine, out *receiptingmodel.ReceiptLine) *receiptingmodel.ReceiptLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptingmodel.ReceiptLine{}
	}
	convert_receiptmodel_ReceiptLine_receiptingmodel_ReceiptLine(arg, out)
	return out
}

func convert_receiptmodel_ReceiptLine_receiptingmodel_ReceiptLine(arg *receiptmodel.ReceiptLine, out *receiptingmodel.ReceiptLine) {
	out.RefID = arg.RefID   // simple assign
	out.Title = arg.Title   // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_receiptmodel_ReceiptLines_receiptingmodel_ReceiptLines(args []*receiptmodel.ReceiptLine) (outs []*receiptingmodel.ReceiptLine) {
	tmps := make([]receiptingmodel.ReceiptLine, len(args))
	outs = make([]*receiptingmodel.ReceiptLine, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptmodel_ReceiptLine_receiptingmodel_ReceiptLine(args[i], &tmps[i])
	}
	return outs
}

func Convert_receiptingmodel_ReceiptLine_receiptmodel_ReceiptLine(arg *receiptingmodel.ReceiptLine, out *receiptmodel.ReceiptLine) *receiptmodel.ReceiptLine {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptmodel.ReceiptLine{}
	}
	convert_receiptingmodel_ReceiptLine_receiptmodel_ReceiptLine(arg, out)
	return out
}

func convert_receiptingmodel_ReceiptLine_receiptmodel_ReceiptLine(arg *receiptingmodel.ReceiptLine, out *receiptmodel.ReceiptLine) {
	out.RefID = arg.RefID   // simple assign
	out.Title = arg.Title   // simple assign
	out.Amount = arg.Amount // simple assign
}

func Convert_receiptingmodel_ReceiptLines_receiptmodel_ReceiptLines(args []*receiptingmodel.ReceiptLine) (outs []*receiptmodel.ReceiptLine) {
	tmps := make([]receiptmodel.ReceiptLine, len(args))
	outs = make([]*receiptmodel.ReceiptLine, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptingmodel_ReceiptLine_receiptmodel_ReceiptLine(args[i], &tmps[i])
	}
	return outs
}

//-- convert etop.vn/backend/com/main/receipting/model.Trader --//

func Convert_receiptmodel_Trader_receiptingmodel_Trader(arg *receiptmodel.Trader, out *receiptingmodel.Trader) *receiptingmodel.Trader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptingmodel.Trader{}
	}
	convert_receiptmodel_Trader_receiptingmodel_Trader(arg, out)
	return out
}

func convert_receiptmodel_Trader_receiptingmodel_Trader(arg *receiptmodel.Trader, out *receiptingmodel.Trader) {
	out.ID = arg.ID             // simple assign
	out.Type = arg.Type         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Phone = arg.Phone       // simple assign
}

func Convert_receiptmodel_Traders_receiptingmodel_Traders(args []*receiptmodel.Trader) (outs []*receiptingmodel.Trader) {
	tmps := make([]receiptingmodel.Trader, len(args))
	outs = make([]*receiptingmodel.Trader, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptmodel_Trader_receiptingmodel_Trader(args[i], &tmps[i])
	}
	return outs
}

func Convert_receiptingmodel_Trader_receiptmodel_Trader(arg *receiptingmodel.Trader, out *receiptmodel.Trader) *receiptmodel.Trader {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &receiptmodel.Trader{}
	}
	convert_receiptingmodel_Trader_receiptmodel_Trader(arg, out)
	return out
}

func convert_receiptingmodel_Trader_receiptmodel_Trader(arg *receiptingmodel.Trader, out *receiptmodel.Trader) {
	out.ID = arg.ID             // simple assign
	out.Type = arg.Type         // simple assign
	out.FullName = arg.FullName // simple assign
	out.Phone = arg.Phone       // simple assign
}

func Convert_receiptingmodel_Traders_receiptmodel_Traders(args []*receiptingmodel.Trader) (outs []*receiptmodel.Trader) {
	tmps := make([]receiptmodel.Trader, len(args))
	outs = make([]*receiptmodel.Trader, len(args))
	for i := range tmps {
		outs[i] = Convert_receiptingmodel_Trader_receiptmodel_Trader(args[i], &tmps[i])
	}
	return outs
}
