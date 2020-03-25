// +build !generator

// Code generated by generator convert. DO NOT EDIT.

package convert

import (
	time "time"

	moneytxmodel "etop.vn/backend/com/main/moneytx/model"
	conversion "etop.vn/backend/pkg/common/conversion"
	moneytransactionshippingmodel "etop.vn/backend/zexp/etl/main/moneytransactionshipping/model"
)

/*
Custom conversions: (none)

Ignored functions: (none)
*/

func RegisterConversions(s *conversion.Scheme) {
	registerConversions(s)
}

func registerConversions(s *conversion.Scheme) {
	s.Register((*moneytransactionshippingmodel.MoneyTransactionShipping)(nil), (*moneytxmodel.MoneyTransactionShipping)(nil), func(arg, out interface{}) error {
		Convert_moneytransactionshippingmodel_MoneyTransactionShipping_moneytxmodel_MoneyTransactionShipping(arg.(*moneytransactionshippingmodel.MoneyTransactionShipping), out.(*moneytxmodel.MoneyTransactionShipping))
		return nil
	})
	s.Register(([]*moneytransactionshippingmodel.MoneyTransactionShipping)(nil), (*[]*moneytxmodel.MoneyTransactionShipping)(nil), func(arg, out interface{}) error {
		out0 := Convert_moneytransactionshippingmodel_MoneyTransactionShippings_moneytxmodel_MoneyTransactionShippings(arg.([]*moneytransactionshippingmodel.MoneyTransactionShipping))
		*out.(*[]*moneytxmodel.MoneyTransactionShipping) = out0
		return nil
	})
	s.Register((*moneytxmodel.MoneyTransactionShipping)(nil), (*moneytransactionshippingmodel.MoneyTransactionShipping)(nil), func(arg, out interface{}) error {
		Convert_moneytxmodel_MoneyTransactionShipping_moneytransactionshippingmodel_MoneyTransactionShipping(arg.(*moneytxmodel.MoneyTransactionShipping), out.(*moneytransactionshippingmodel.MoneyTransactionShipping))
		return nil
	})
	s.Register(([]*moneytxmodel.MoneyTransactionShipping)(nil), (*[]*moneytransactionshippingmodel.MoneyTransactionShipping)(nil), func(arg, out interface{}) error {
		out0 := Convert_moneytxmodel_MoneyTransactionShippings_moneytransactionshippingmodel_MoneyTransactionShippings(arg.([]*moneytxmodel.MoneyTransactionShipping))
		*out.(*[]*moneytransactionshippingmodel.MoneyTransactionShipping) = out0
		return nil
	})
}

//-- convert etop.vn/backend/com/main/moneytx/model.MoneyTransactionShipping --//

func Convert_moneytransactionshippingmodel_MoneyTransactionShipping_moneytxmodel_MoneyTransactionShipping(arg *moneytransactionshippingmodel.MoneyTransactionShipping, out *moneytxmodel.MoneyTransactionShipping) *moneytxmodel.MoneyTransactionShipping {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &moneytxmodel.MoneyTransactionShipping{}
	}
	convert_moneytransactionshippingmodel_MoneyTransactionShipping_moneytxmodel_MoneyTransactionShipping(arg, out)
	return out
}

func convert_moneytransactionshippingmodel_MoneyTransactionShipping_moneytxmodel_MoneyTransactionShipping(arg *moneytransactionshippingmodel.MoneyTransactionShipping, out *moneytxmodel.MoneyTransactionShipping) {
	out.ID = arg.ID                            // simple assign
	out.ShopID = arg.ShopID                    // simple assign
	out.CreatedAt = arg.CreatedAt              // simple assign
	out.UpdatedAt = arg.UpdatedAt              // simple assign
	out.ClosedAt = arg.ClosedAt                // simple assign
	out.Status = arg.Status                    // simple assign
	out.TotalCOD = arg.TotalCOD                // simple assign
	out.TotalAmount = arg.TotalAmount          // simple assign
	out.TotalOrders = arg.TotalOrders          // simple assign
	out.Code = arg.Code                        // simple assign
	out.MoneyTransactionShippingExternalID = 0 // zero value
	out.MoneyTransactionShippingEtopID = 0     // zero value
	out.Provider = arg.Provider                // simple assign
	out.ConfirmedAt = arg.ConfirmedAt          // simple assign
	out.EtopTransferedAt = time.Time{}         // zero value
	out.BankAccount = arg.BankAccount          // simple assign
	out.Note = arg.Note                        // simple assign
	out.InvoiceNumber = arg.InvoiceNumber      // simple assign
	out.Type = arg.Type                        // simple assign
	out.Rid = arg.Rid                          // simple assign
}

func Convert_moneytransactionshippingmodel_MoneyTransactionShippings_moneytxmodel_MoneyTransactionShippings(args []*moneytransactionshippingmodel.MoneyTransactionShipping) (outs []*moneytxmodel.MoneyTransactionShipping) {
	tmps := make([]moneytxmodel.MoneyTransactionShipping, len(args))
	outs = make([]*moneytxmodel.MoneyTransactionShipping, len(args))
	for i := range tmps {
		outs[i] = Convert_moneytransactionshippingmodel_MoneyTransactionShipping_moneytxmodel_MoneyTransactionShipping(args[i], &tmps[i])
	}
	return outs
}

func Convert_moneytxmodel_MoneyTransactionShipping_moneytransactionshippingmodel_MoneyTransactionShipping(arg *moneytxmodel.MoneyTransactionShipping, out *moneytransactionshippingmodel.MoneyTransactionShipping) *moneytransactionshippingmodel.MoneyTransactionShipping {
	if arg == nil {
		return nil
	}
	if out == nil {
		out = &moneytransactionshippingmodel.MoneyTransactionShipping{}
	}
	convert_moneytxmodel_MoneyTransactionShipping_moneytransactionshippingmodel_MoneyTransactionShipping(arg, out)
	return out
}

func convert_moneytxmodel_MoneyTransactionShipping_moneytransactionshippingmodel_MoneyTransactionShipping(arg *moneytxmodel.MoneyTransactionShipping, out *moneytransactionshippingmodel.MoneyTransactionShipping) {
	out.ID = arg.ID                       // simple assign
	out.ShopID = arg.ShopID               // simple assign
	out.CreatedAt = arg.CreatedAt         // simple assign
	out.UpdatedAt = arg.UpdatedAt         // simple assign
	out.ClosedAt = arg.ClosedAt           // simple assign
	out.Status = arg.Status               // simple assign
	out.TotalCOD = arg.TotalCOD           // simple assign
	out.TotalAmount = arg.TotalAmount     // simple assign
	out.TotalOrders = arg.TotalOrders     // simple assign
	out.Code = arg.Code                   // simple assign
	out.Provider = arg.Provider           // simple assign
	out.ConfirmedAt = arg.ConfirmedAt     // simple assign
	out.BankAccount = arg.BankAccount     // simple assign
	out.Note = arg.Note                   // simple assign
	out.InvoiceNumber = arg.InvoiceNumber // simple assign
	out.Type = arg.Type                   // simple assign
	out.Rid = arg.Rid                     // simple assign
}

func Convert_moneytxmodel_MoneyTransactionShippings_moneytransactionshippingmodel_MoneyTransactionShippings(args []*moneytxmodel.MoneyTransactionShipping) (outs []*moneytransactionshippingmodel.MoneyTransactionShipping) {
	tmps := make([]moneytransactionshippingmodel.MoneyTransactionShipping, len(args))
	outs = make([]*moneytransactionshippingmodel.MoneyTransactionShipping, len(args))
	for i := range tmps {
		outs[i] = Convert_moneytxmodel_MoneyTransactionShipping_moneytransactionshippingmodel_MoneyTransactionShipping(args[i], &tmps[i])
	}
	return outs
}
