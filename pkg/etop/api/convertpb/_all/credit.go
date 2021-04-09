package _all

import (
	"o.o/api/main/credit"
	"o.o/api/top/int/etop"
	creditmodel "o.o/backend/com/main/credit/model"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

func PbCreditExtended(item *creditmodel.CreditExtended) *etop.Credit {
	if item == nil {
		return nil
	}

	return &etop.Credit{
		ID:        item.ID,
		Amount:    item.Amount,
		ShopID:    item.ShopID,
		Type:      item.Type,
		Shop:      convertpb.PbShop(item.Shop),
		CreatedAt: cmapi.PbTime(item.CreatedAt),
		UpdatedAt: cmapi.PbTime(item.UpdatedAt),
		PaidAt:    cmapi.PbTime(item.PaidAt),
		Status:    item.Status,
	}
}

func PbCreditExtendeds(items []*creditmodel.CreditExtended) []*etop.Credit {
	result := make([]*etop.Credit, len(items))
	for i, item := range items {
		result[i] = PbCreditExtended(item)
	}
	return result
}

func Convert_core_CreditExtended_to_api_Credit(item *credit.CreditExtended) *etop.Credit {
	if item == nil {
		return nil
	}

	return &etop.Credit{
		ID:        item.ID,
		Amount:    item.Amount,
		ShopID:    item.ShopID,
		Type:      item.Type,
		Shop:      Convert_core_Shop_To_api_Shop(item.Shop, nil),
		CreatedAt: cmapi.PbTime(item.CreatedAt),
		UpdatedAt: cmapi.PbTime(item.UpdatedAt),
		PaidAt:    cmapi.PbTime(item.PaidAt),
		Status:    item.Status,
		Classify:  item.Classify,
	}
}

func Convert_core_CreditExtendeds_to_api_Credits(items []*credit.CreditExtended) []*etop.Credit {
	result := make([]*etop.Credit, len(items))
	for i, item := range items {
		result[i] = Convert_core_CreditExtended_to_api_Credit(item)
	}
	return result
}

func Convert_core_Credit_to_api_Credit(item *credit.Credit) *etop.Credit {
	if item == nil {
		return nil
	}
	return &etop.Credit{
		ID:        item.ID,
		Amount:    item.Amount,
		ShopID:    item.ShopID,
		Type:      item.Type,
		CreatedAt: cmapi.PbTime(item.CreatedAt),
		UpdatedAt: cmapi.PbTime(item.UpdatedAt),
		PaidAt:    cmapi.PbTime(item.PaidAt),
		Status:    item.Status,
		Classify:  item.Classify,
	}
}

func Convert_core_Credits_to_api_Credits(items []*credit.Credit) []*etop.Credit {
	result := make([]*etop.Credit, len(items))
	for i, item := range items {
		result[i] = Convert_core_Credit_to_api_Credit(item)
	}
	return result
}
