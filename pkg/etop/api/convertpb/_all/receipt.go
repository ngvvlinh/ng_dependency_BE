package _all

import (
	"o.o/api/main/receipting"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_api_ReceiptLine_To_core_ReceiptLine(in *shop.ReceiptLine) *receipting.ReceiptLine {
	if in == nil {
		return nil
	}
	return &receipting.ReceiptLine{
		RefID:  in.RefId,
		Title:  in.Title,
		Amount: in.Amount,
	}
}

func Convert_api_ReceiptLines_To_core_ReceiptLines(in []*shop.ReceiptLine) []*receipting.ReceiptLine {
	out := make([]*receipting.ReceiptLine, len(in))
	for i := range in {
		out[i] = Convert_api_ReceiptLine_To_core_ReceiptLine(in[i])
	}

	return out
}

func PbReceiptLine(m *receipting.ReceiptLine) *shop.ReceiptLine {
	if m == nil {
		return nil
	}
	return &shop.ReceiptLine{
		RefId:  m.RefID,
		Title:  m.Title,
		Amount: m.Amount,
	}
}

func PbReceiptLines(ms []*receipting.ReceiptLine) []*shop.ReceiptLine {
	res := make([]*shop.ReceiptLine, len(ms))
	for i, m := range ms {
		res[i] = PbReceiptLine(m)
	}
	return res
}

func PbReceipt(m *receipting.Receipt) *shop.Receipt {
	return &shop.Receipt{
		Id:           m.ID,
		ShopId:       m.ShopID,
		TraderId:     m.TraderID,
		Code:         m.Code,
		Title:        m.Title,
		Type:         m.Type,
		Description:  m.Description,
		Amount:       m.Amount,
		LedgerId:     m.LedgerID,
		RefType:      m.RefType,
		CancelReason: m.CancelReason,
		Lines:        PbReceiptLines(m.Lines),
		Trader:       PbTrader(m.Trader),
		Status:       m.Status,
		CreatedBy:    m.CreatedBy,
		CreatedType:  m.Mode,
		Mode:         m.Mode,
		PaidAt:       cmapi.PbTime(m.PaidAt),
		ConfirmedAt:  cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:  cmapi.PbTime(m.CancelledAt),
		CreatedAt:    cmapi.PbTime(m.CreatedAt),
		UpdatedAt:    cmapi.PbTime(m.UpdatedAt),
		Note:         m.Note,
	}
}

func PbReceipts(ms []*receipting.Receipt) []*shop.Receipt {
	res := make([]*shop.Receipt, len(ms))
	for i, m := range ms {
		res[i] = PbReceipt(m)
	}
	return res
}

func PbTrader(m *receipting.Trader) *shop.Trader {
	if m == nil {
		return nil
	}
	return &shop.Trader{
		Id:       m.ID,
		Type:     m.Type,
		FullName: m.FullName,
		Phone:    m.Phone,
		Deleted:  false,
	}
}

func PbTraders(ms []*receipting.Trader) []*shop.Trader {
	res := make([]*shop.Trader, len(ms))
	for i, m := range ms {
		res[i] = PbTrader(m)
	}
	return res
}
