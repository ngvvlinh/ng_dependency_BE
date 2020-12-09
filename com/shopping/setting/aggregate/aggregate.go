package aggregate

import (
	"context"

	addressing "o.o/api/main/address"
	"o.o/api/shopping/setting"
	"o.o/api/top/types/etc/address_type"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/setting/convert"
	"o.o/backend/com/shopping/setting/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/capi/dot"
)

var _ setting.Aggregate = &ShopSettingAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type ShopSettingAggregate struct {
	store       sqlstore.ShopSettingStoreFactory
	addressAggr addressing.CommandBus
}

func NewShopSettingAggregate(
	db com.MainDB, addressA addressing.CommandBus,
) *ShopSettingAggregate {
	return &ShopSettingAggregate{
		store:       sqlstore.NewShopSettingStore(db),
		addressAggr: addressA,
	}
}

func ShopSettingAggregateMessageBus(a *ShopSettingAggregate) setting.CommandBus {
	b := bus.New()
	return setting.NewAggregateHandler(a).RegisterHandlers(b)
}

func (s *ShopSettingAggregate) UpdateShopSetting(
	ctx context.Context, args *setting.UpdateShopSettingArgs,
) (*setting.ShopSetting, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}

	shopSetting := new(setting.ShopSetting)
	if err := scheme.Convert(args, shopSetting); err != nil {
		return nil, err
	}

	oldShopSetting, err := s.store(ctx).ShopID(args.ShopID).GetShopSetting()
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		// create new shopSetting
		if args.ReturnAddress != nil {
			returnAddress, err := s.createReturnAddress(ctx, args.ReturnAddress, args.ShopID)
			if err != nil {
				return nil, err
			}

			shopSetting.ReturnAddressID = returnAddress.ID
		}

		if err := s.store(ctx).CreateShopSetting(shopSetting); err != nil {
			return nil, err
		}
	case cm.NoError:
		// update shopSetting
		if args.ReturnAddress != nil {
			if oldShopSetting.ReturnAddressID == 0 {
				returnAddress, err := s.createReturnAddress(ctx, args.ReturnAddress, args.ShopID)
				if err != nil {
					return nil, err
				}

				shopSetting.ReturnAddressID = returnAddress.ID
			} else {
				returnAddress, err := s.updateReturnAddress(ctx, oldShopSetting.ReturnAddressID, args.ReturnAddress)
				if err != nil {
					return nil, err
				}

				shopSetting.ReturnAddressID = returnAddress.ID
			}
		}

		if err := s.store(ctx).UpdateShopSetting(shopSetting); err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return shopSetting, nil
}

func (s *ShopSettingAggregate) updateReturnAddress(
	ctx context.Context, ID dot.ID, returnAddress *addressing.Address,
) (*addressing.Address, error) {
	updateReturnAddressCmd := &addressing.UpdateAddressCommand{
		ID:           ID,
		Province:     returnAddress.Province,
		ProvinceCode: returnAddress.ProvinceCode,
		District:     returnAddress.District,
		DistrictCode: returnAddress.DistrictCode,
		Ward:         returnAddress.Ward,
		WardCode:     returnAddress.WardCode,
		Address1:     returnAddress.Address1,
		Address2:     returnAddress.Address2,
		Zip:          returnAddress.Zip,
		Company:      returnAddress.Company,
		City:         returnAddress.City,
		Country:      returnAddress.Country,
		FullName:     returnAddress.FullName,
		FirstName:    returnAddress.FirstName,
		LastName:     returnAddress.LastName,
		Phone:        returnAddress.Phone,
		Email:        returnAddress.Email,
		Position:     returnAddress.Position,
		Notes:        returnAddress.Notes,
		Type:         address_type.Warehouse,
		Coordinates:  returnAddress.Coordinates,
	}

	if err := s.addressAggr.Dispatch(ctx, updateReturnAddressCmd); err != nil {
		return nil, err
	}
	return updateReturnAddressCmd.Result, nil
}

func (s *ShopSettingAggregate) createReturnAddress(
	ctx context.Context, returnAddress *addressing.Address, shopID dot.ID,
) (*addressing.Address, error) {
	createReturnAddressCmd := &addressing.CreateAddressCommand{
		Province:     returnAddress.Province,
		ProvinceCode: returnAddress.ProvinceCode,
		District:     returnAddress.District,
		DistrictCode: returnAddress.DistrictCode,
		Ward:         returnAddress.Ward,
		WardCode:     returnAddress.WardCode,
		Address1:     returnAddress.Address1,
		Address2:     returnAddress.Address2,
		Zip:          returnAddress.Zip,
		Company:      returnAddress.Company,
		City:         returnAddress.City,
		Country:      returnAddress.Country,
		FullName:     returnAddress.FullName,
		FirstName:    returnAddress.FirstName,
		LastName:     returnAddress.LastName,
		Phone:        returnAddress.Phone,
		Email:        returnAddress.Email,
		Position:     returnAddress.Position,
		Type:         address_type.Warehouse,
		AccountID:    shopID,
		Notes:        returnAddress.Notes,
		Coordinates:  returnAddress.Coordinates,
	}
	if err := s.addressAggr.Dispatch(ctx, createReturnAddressCmd); err != nil {
		return nil, err
	}
	return createReturnAddressCmd.Result, nil
}
