package _all

import (
	"o.o/api/main/moneytx"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/account_tag"
	txmodel "o.o/backend/com/main/moneytx/model"
	"o.o/backend/com/main/moneytx/txmodely"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
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
		BankAccount:    convertpb.Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
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
		BankAccount:                        convertpb.Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
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
		BankAccount:           convertpb.Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
	}
}

func PbMoneyTxShippingEtops(items []*moneytx.MoneyTransactionShippingEtop) []*types.MoneyTransactionShippingEtop {
	result := make([]*types.MoneyTransactionShippingEtop, len(items))
	for i, item := range items {
		result[i] = PbMoneyTxShippingEtop(item)
	}
	return result
}

func PbMoneyTransactionExtendeds(items []*txmodely.MoneyTransactionExtended) []*types.MoneyTransaction {
	result := make([]*types.MoneyTransaction, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingEtopExtended(m *txmodely.MoneyTransactionShippingEtopExtended) *types.MoneyTransactionShippingEtop {
	if m == nil {
		return nil
	}
	return &types.MoneyTransactionShippingEtop{
		Id:                m.ID,
		Code:              m.Code,
		TotalCod:          m.TotalCOD,
		TotalOrders:       m.TotalOrders,
		TotalAmount:       m.TotalAmount,
		TotalFee:          m.TotalFee,
		Status:            m.Status,
		MoneyTransactions: PbMoneyTransactionExtendeds(m.MoneyTransactions),
		CreatedAt:         cmapi.PbTime(m.CreatedAt),
		UpdatedAt:         cmapi.PbTime(m.UpdatedAt),
		ConfirmedAt:       cmapi.PbTime(m.ConfirmedAt),
		Note:              m.Note,
		InvoiceNumber:     m.InvoiceNumber,
		BankAccount:       convertpb.PbBankAccount(m.BankAccount),
	}
}

func PbMoneyTransactionShippingEtopExtendeds(items []*txmodely.MoneyTransactionShippingEtopExtended) []*types.MoneyTransactionShippingEtop {
	result := make([]*types.MoneyTransactionShippingEtop, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingEtopExtended(item)
	}
	return result
}

func PbMoneyTransactionExtended(m *txmodely.MoneyTransactionExtended) *types.MoneyTransaction {
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
	}
}

func PbMoneyTransaction(m *txmodel.MoneyTransactionShipping) *types.MoneyTransaction {
	if m == nil {
		return nil
	}
	return &types.MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             m.Status,
		TotalCod:                           m.TotalCOD,
		TotalOrders:                        m.TotalOrders,
		Code:                               m.Code,
		Provider:                           m.Provider,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:                        cmapi.PbTime(m.ConfirmedAt),
		EtopTransferedAt:                   cmapi.PbTime(m.EtopTransferedAt),
	}
}

func PbMoneyTransactionShippingExternalExtended(m *txmodel.MoneyTransactionShippingExternalExtended) *types.MoneyTransactionShippingExternal {
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
		CreatedAt:      cmapi.PbTime(m.CreatedAt),
		UpdatedAt:      cmapi.PbTime(m.UpdatedAt),
		ExternalPaidAt: cmapi.PbTime(m.ExternalPaidAt),
		Note:           m.Note,
		InvoiceNumber:  m.InvoiceNumber,
		BankAccount:    convertpb.PbBankAccount(m.BankAccount),
		Lines:          PbMoneyTransactionShippingExternalLineExtendeds(m.Lines),
	}

	return res
}

func PbMoneyTransactionShippingExternalExtendeds(items []*txmodel.MoneyTransactionShippingExternalExtended) []*types.MoneyTransactionShippingExternal {
	result := make([]*types.MoneyTransactionShippingExternal, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalLineExtended(m *txmodel.MoneyTransactionShippingExternalLineExtended) *types.MoneyTransactionShippingExternalLine {
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
		ImportError:                        cmapi.PbCustomError(m.ImportError),
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ExternalCreatedAt:                  cmapi.PbTime(m.ExternalCreatedAt),
		ExternalClosedAt:                   cmapi.PbTime(m.ExternalClosedAt),
	}
	if m.Fulfillment != nil && m.Fulfillment.ID != 0 {
		res.Fulfillment = convertpb.PbFulfillment(m.Fulfillment, account_tag.TagEtop, m.Shop, m.Order)
	}
	return res
}

func PbMoneyTransactionShippingExternalLineExtendeds(items []*txmodel.MoneyTransactionShippingExternalLineExtended) []*types.MoneyTransactionShippingExternalLine {
	result := make([]*types.MoneyTransactionShippingExternalLine, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalLineExtended(item)
	}
	return result
}
