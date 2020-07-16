package convertpb

import (
	"o.o/api/main/moneytx"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbMoneyTxShippingExternalsFtLine(items []*moneytx.MoneyTransactionShippingExternalFtLine) []*types.MoneyTransactionShippingExternal {
	result := make([]*types.MoneyTransactionShippingExternal, len(items))
	for i, item := range items {
		result[i] = PbMoneyTxShippingExternalFtLine(item)
	}
	return result
}

func PbMoneyTxShippingExternalFtLine(m *moneytx.MoneyTransactionShippingExternalFtLine) *types.MoneyTransactionShippingExternal {
	if m == nil {
		return nil
	}
	res := &types.MoneyTransactionShippingExternal{
		Id:             m.ID,
		Code:           m.Code,
		TotalCod:       m.TotalCOD,
		TotalOrders:    m.TotalOrders,
		Status:         m.Status,
		Provider:       m.Provider.String(),
		Lines:          PbMoneyTxShippingExternalLines(m.Lines),
		CreatedAt:      cmapi.PbTime(m.CreatedAt),
		UpdatedAt:      cmapi.PbTime(m.UpdatedAt),
		ExternalPaidAt: cmapi.PbTime(m.ExternalPaidAt),
		Note:           m.Note,
		InvoiceNumber:  m.InvoiceNumber,
		BankAccount:    Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
		ConnectionID:   m.ConnectionID,
	}
	return res
}

func PbMoneyTxShippingExternalLines(items []*moneytx.MoneyTransactionShippingExternalLine) []*types.MoneyTransactionShippingExternalLine {
	result := make([]*types.MoneyTransactionShippingExternalLine, len(items))
	for i, item := range items {
		result[i] = PbMoneyTxShippingExternalLine(item)
	}
	return result
}

func PbMoneyTxShippingExternalLine(m *moneytx.MoneyTransactionShippingExternalLine) *types.MoneyTransactionShippingExternalLine {
	if m == nil {
		return nil
	}
	res := &types.MoneyTransactionShippingExternalLine{
		Id:                                 m.ID,
		ExternalCode:                       m.ExternalCode,
		ExternalCustomer:                   m.ExternalCustomer,
		ExternalAddress:                    m.ExternalAddress,
		ExternalTotalCod:                   m.ExternalTotalCOD,
		ExternalTotalShippingFee:           m.ExternalTotalShippingFee,
		EtopFulfillmentId:                  m.EtopFulfillmentID,
		EtopFulfillmentIdRaw:               m.EtopFulfillmentIDRaw,
		Note:                               m.Note,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		ImportError:                        cmapi.PbMetaError(m.ImportError),
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ExternalCreatedAt:                  cmapi.PbTime(m.ExternalCreatedAt),
		ExternalClosedAt:                   cmapi.PbTime(m.ExternalClosedAt),
	}
	return res
}

func PbMoneyTxShipping(m *moneytx.MoneyTransactionShipping) *types.MoneyTransaction {
	if m == nil {
		return nil
	}
	return &types.MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             m.Status,
		TotalCod:                           m.TotalCOD,
		TotalOrders:                        m.TotalOrders,
		TotalAmount:                        m.TotalAmount,
		Code:                               m.Code,
		Provider:                           m.Provider,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		MoneyTransactionShippingEtopId:     m.MoneyTransactionShippingEtopID,
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:                        cmapi.PbTime(m.ConfirmedAt),
		EtopTransferedAt:                   cmapi.PbTime(m.EtopTransferedAt),
		Note:                               m.Note,
		BankAccount:                        Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
		InvoiceNumber:                      m.InvoiceNumber,
	}
}

func PbMoneyTxShippings(items []*moneytx.MoneyTransactionShipping) []*types.MoneyTransaction {
	result := make([]*types.MoneyTransaction, len(items))
	for i, item := range items {
		result[i] = PbMoneyTxShipping(item)
	}
	return result
}

func PbMoneyTxShippingEtop(m *moneytx.MoneyTransactionShippingEtop) *types.MoneyTransactionShippingEtop {
	if m == nil {
		return nil
	}
	return &types.MoneyTransactionShippingEtop{
		Id:                    m.ID,
		Code:                  m.Code,
		TotalCod:              m.TotalCOD,
		TotalOrders:           m.TotalOrders,
		TotalAmount:           m.TotalAmount,
		TotalFee:              m.TotalFee,
		TotalMoneyTxShippings: m.TotalMoneyTransaction,
		Status:                m.Status,
		CreatedAt:             cmapi.PbTime(m.CreatedAt),
		UpdatedAt:             cmapi.PbTime(m.UpdatedAt),
		ConfirmedAt:           cmapi.PbTime(m.ConfirmedAt),
		Note:                  m.Note,
		InvoiceNumber:         m.InvoiceNumber,
		BankAccount:           Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
	}
}

func PbMoneyTxShippingEtops(items []*moneytx.MoneyTransactionShippingEtop) []*types.MoneyTransactionShippingEtop {
	result := make([]*types.MoneyTransactionShippingEtop, len(items))
	for i, item := range items {
		result[i] = PbMoneyTxShippingEtop(item)
	}
	return result
}
