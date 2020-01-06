package convertpb

import (
	"etop.vn/api/main/moneytx"
	"etop.vn/api/top/int/types"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/model"
)

func PbMoneyTxShippingExternalExtended(m *moneytx.MoneyTransactionShippingExternalExtended) *types.MoneyTransactionShippingExternal {
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
		Lines:          PbMoneyTxShippingExternalLineExtendeds(m.Lines),
		CreatedAt:      cmapi.PbTime(m.CreatedAt),
		UpdatedAt:      cmapi.PbTime(m.UpdatedAt),
		ExternalPaidAt: cmapi.PbTime(m.ExternalPaidAt),
		Note:           m.Note,
		InvoiceNumber:  m.InvoiceNumber,
		BankAccount:    Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
	}
	return res
}

func PbMoneyTxShippingExternalLineExtendeds(items []*moneytx.MoneyTransactionShippingExternalLineExtended) []*types.MoneyTransactionShippingExternalLine {
	result := make([]*types.MoneyTransactionShippingExternalLine, len(items))
	for i, item := range items {
		result[i] = PbMoneyTxShippingExternalLineExtended(item)
	}
	return result
}

func PbMoneyTxShippingExternalLineExtended(m *moneytx.MoneyTransactionShippingExternalLineExtended) *types.MoneyTransactionShippingExternalLine {
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
	if m.Fulfillment != nil && m.Fulfillment.ID != 0 {
		res.Fulfillment = Convert_core_Fulfillment_To_api_Fulfillment(m.Fulfillment, model.TagEtop, m.Shop, m.Order)
	}
	return res
}
