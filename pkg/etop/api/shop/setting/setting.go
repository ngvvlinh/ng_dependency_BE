package setting

import (
	"context"

	addressing "o.o/api/main/address"
	"o.o/api/shopping/setting"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type SettingService struct {
	session.Session

	SettingQuery setting.QueryBus
	SettingAggr  setting.CommandBus

	AddressQ addressing.QueryBus
}

func (s *SettingService) Close() api.SettingService { res := *s; return &res }

func (s *SettingService) GetSetting(
	ctx context.Context, empty *pbcm.Empty,
) (*api.ShopSetting, error) {
	getShopSettingQuery := &setting.GetShopSettingQuery{
		ShopID: s.SS.Shop().ID,
	}
	if err := s.SettingQuery.Dispatch(ctx, getShopSettingQuery); err != nil {
		return nil, err
	}

	shopSetting := getShopSettingQuery.Result
	if err := s.GetAddress(ctx, shopSetting); err != nil {
		return nil, err
	}

	return convertpb.Convert_core_Setting_to_api_Setting(shopSetting), nil
}

func (s *SettingService) UpdateSetting(
	ctx context.Context, req *api.UpdateSettingRequest,
) (*api.ShopSetting, error) {
	updateShopSettingCmd := &setting.UpdateShopSettingCommand{
		ShopID:          s.SS.Shop().ID,
		ReturnAddress:   convertpb.Convert_api_EtopAddress_To_core_Address(req.ReturnAddress),
		PaymentTypeID:   req.PaymentTypeID,
		TryOn:           req.TryOn,
		ShippingNote:    req.ShippingNote,
		Weight:          req.Weight,
		HideAllComments: req.HideAllComments,
	}
	if err := s.SettingAggr.Dispatch(ctx, updateShopSettingCmd); err != nil {
		return nil, err
	}

	shopSetting := updateShopSettingCmd.Result
	if err := s.GetAddress(ctx, shopSetting); err != nil {
		return nil, err
	}

	return convertpb.Convert_core_Setting_to_api_Setting(shopSetting), nil
}

func (s *SettingService) GetAddress(ctx context.Context, shopSetting *setting.ShopSetting) error {
	if shopSetting.ReturnAddressID == 0 {
		return nil
	}

	getAddressQuery := &addressing.GetAddressByIDQuery{
		ID: shopSetting.ReturnAddressID,
	}
	if err := s.AddressQ.Dispatch(ctx, getAddressQuery); err != nil {
		return err
	}

	shopSetting.ReturnAddress = getAddressQuery.Result
	return nil
}
