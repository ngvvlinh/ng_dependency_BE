package _all

import (
	"o.o/api/main/identity"
	"o.o/api/main/moneytx"
	"o.o/api/top/int/etop"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
)

func Convert_core_Shop_To_api_Shop(in *identity.Shop, shopCount *moneytx.ShopFtMoneyTxShippingCount) *etop.Shop {
	if in == nil {
		return nil
	}
	res := &etop.Shop{
		Id:                      in.ID,
		Name:                    in.Name,
		Status:                  in.Status,
		IsTest:                  in.IsTest == 1,
		Phone:                   in.Phone,
		AutoCreateFfm:           in.AutoCreateFFM,
		WebsiteUrl:              in.WebsiteURL,
		ImageUrl:                in.ImageURL,
		Email:                   in.Email,
		ShipToAddressId:         in.ShipToAddressID,
		ShipFromAddressId:       in.ShipFromAddressID,
		OwnerId:                 in.OwnerID,
		Code:                    in.Code,
		BankAccount:             convertpb.Convert_core_BankAccount_To_api_BankAccount(in.BankAccount),
		TryOn:                   in.TryOn,
		MoneyTransactionRrule:   in.MoneyTransactionRRule,
		IsPriorMoneyTransaction: in.IsPriorMoneyTransaction,
	}
	if shopCount != nil && shopCount.ShopID == in.ID {
		res.MoneyTransactionCount = shopCount.MoneyTxShippingCount
	}
	return res
}

func Convert_core_Shops_To_api_Shops(items []*identity.Shop, shopFtMoneyTxShippingCounts []*moneytx.ShopFtMoneyTxShippingCount) []*etop.Shop {
	result := make([]*etop.Shop, len(items))

	mapShopFtMoneyTxShippingCounts := make(map[dot.ID]*moneytx.ShopFtMoneyTxShippingCount)
	if shopFtMoneyTxShippingCounts != nil {
		for _, shopCount := range shopFtMoneyTxShippingCounts {
			mapShopFtMoneyTxShippingCounts[shopCount.ShopID] = shopCount
		}
	}

	for i, item := range items {
		result[i] = Convert_core_Shop_To_api_Shop(item, mapShopFtMoneyTxShippingCounts[item.ID])
	}
	return result
}

func Convert_core_ShopExtended_To_api_ShopExtended(m *identity.ShopExtended, shopCount *moneytx.ShopFtMoneyTxShippingCount) *etop.Shop {
	if m == nil {
		return nil
	}
	res := &etop.Shop{
		Id:                            m.ID,
		InventoryOverstock:            m.InventoryOverstock.Apply(true),
		Name:                          m.Name,
		Status:                        m.Status,
		Address:                       convertpb.Convert_core_Address_To_api_Address(m.Address),
		Phone:                         m.Phone,
		BankAccount:                   convertpb.Convert_core_BankAccount_To_api_BankAccount(m.BankAccount),
		WebsiteUrl:                    m.WebsiteURL,
		ImageUrl:                      m.ImageURL,
		Email:                         m.Email,
		ShipToAddressId:               m.ShipToAddressID,
		ShipFromAddressId:             m.ShipFromAddressID,
		AutoCreateFfm:                 m.AutoCreateFFM,
		TryOn:                         m.TryOn,
		GhnNoteCode:                   m.GhnNoteCode,
		OwnerId:                       m.OwnerID,
		User:                          convertpb.Convert_core_User_To_api_User(m.User),
		CompanyInfo:                   convertpb.Convert_core_CompanyInfo_To_api_CompanyInfo(m.CompanyInfo),
		MoneyTransactionRrule:         m.MoneyTransactionRRule,
		SurveyInfo:                    convertpb.Convert_core_SurveyInfos_To_api_SurveyInfors(m.SurveyInfo),
		ShippingServiceSelectStrategy: convertpb.Convert_core_ShippingServiceSelectStrategy_To_api_ShippingServiceSelectStrategy(m.ShippingServiceSelectStrategy),
		Code:                          m.Code,
		CreatedAt:                     cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                     cmapi.PbTime(m.UpdatedAt),
		IsPriorMoneyTransaction:       m.IsPriorMoneyTransaction,

		// deprecated: 2018.07.24+14
		ProductSourceId: m.ID,
	}
	if shopCount != nil && shopCount.ShopID == m.ID {
		res.MoneyTransactionCount = shopCount.MoneyTxShippingCount
	}
	return res
}

func Convert_core_ShopExtendeds_To_api_ShopExtendeds(items []*identity.ShopExtended, shopFtMoneyTxShippingCounts []*moneytx.ShopFtMoneyTxShippingCount) []*etop.Shop {
	mapShopFtMoneyTxShippingCounts := make(map[dot.ID]*moneytx.ShopFtMoneyTxShippingCount)
	if shopFtMoneyTxShippingCounts != nil {
		for _, shopCount := range shopFtMoneyTxShippingCounts {
			mapShopFtMoneyTxShippingCounts[shopCount.ShopID] = shopCount
		}
	}
	result := make([]*etop.Shop, len(items))
	for i, item := range items {
		result[i] = Convert_core_ShopExtended_To_api_ShopExtended(item, mapShopFtMoneyTxShippingCounts[item.ID])
	}
	return result
}
